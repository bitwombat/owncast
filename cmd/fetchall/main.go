package main

import (
	"log"

	storage "github.com/bitwombat/listen-later/storage/gcloud"
	video "github.com/bitwombat/listen-later/video/youtube"
)

func main() {
	// First, write MP3 to storage from the YouTube playlist
	pl, err := video.GetPlaylist("PLQaHRuFmgISeszmitcke42BAYv3ngK8aL")
	if err != nil {
		log.Println(err)
	}

	stg := &storage.Gcloud{}

	for _, e := range pl.Entries {
		err = video.DownloadMP3ToStorage(e.Id, stg)
		if err != nil {
			log.Println(err)
		}
	}

	// Next, create the RSS feed from storage, and write it to storage
	/* 	err = rss.RssFromStorage(stg)
	   	if err != nil {
	   		log.Println(err)
	   	} */
}
