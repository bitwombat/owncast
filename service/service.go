package service

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	audio "github.com/bitwombat/listen-later/audio/podcast"
	storage "github.com/bitwombat/listen-later/storage/gcloud"
	video "github.com/bitwombat/listen-later/video/youtube"
)

func handleError(w http.ResponseWriter, msg string, err error) {
	log.Printf("ERROR "+msg+": %v", err)
	fmt.Fprintln(w, "Error. See log.")
}

func UpdateHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	plId := r.URL.Path[1:]
	if plId == "" {
		fmt.Fprintf(w, "No list ID")
		return
	}

	log.Println("INFO: Starting update process")
	fmt.Fprintln(w, "Starting update process")

	// First, grab the playlist off YouTube
	pl, err := video.GetPlaylist(plId)
	if err != nil {
		handleError(w, "getting playlist", err)
		return
	}

	stg := &storage.Gcloud{}

	// Then, figure out the entries that are new - ie. in the playlist but not in storage
	newEntries, err := pl.NotInStorage(stg)
	if err != nil {
		handleError(w, "calculating download", err)
		fmt.Fprintln(w, "Error. See log.")
		return
	}

	// Then, download those
	for i, entry := range newEntries {
		log.Printf("Downloading %d of %d", i, len(pl.Entries))
		err := video.DownloadMP3ToStorage(entry, stg) // entries.DownloadToStorage(downloader, stg)
		if err != nil {
			handleError(w, "downloading", err)
			fmt.Fprintln(w, "Error. See log.")
			return
		}
	}

	// Next, figure out the entries that are to be deleted - ie. in storage but not in the playlist
	deleted, err := pl.NotInPlaylist(stg)
	if err != nil {
		handleError(w, "calculating files to delete", err)
		fmt.Fprintln(w, "Error. See log.")
		return
	}

	// Then, delete those
	err = stg.DeleteList(deleted)
	if err != nil {
		handleError(w, "deleting from storage", err)
		fmt.Fprintln(w, "Error. See log.")
		return
	}

	// Finally, create the RSS feed from the playlist
	feed, err := audio.RSSFromPlaylist(pl)
	if err != nil {
		handleError(w, "creating RSS feed", err)
		fmt.Fprintln(w, "Error. See log.")
		return
	}

	// and write it to storage
	feedAsReader := strings.NewReader(feed)
	stg.Put(feedAsReader, "feed.rss")

	log.Println("INFO: Success.")
	fmt.Fprintln(w, "Success.")
}
