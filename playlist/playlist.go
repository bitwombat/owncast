package playlist

import (
	"fmt"

	"github.com/bitwombat/listen-later/storage"
)

type Playlist struct {
	Title   string
	Entries []Entry
}

type Entry struct {
	Id          string
	Title       string
	Upload_Date string
	Webpage_Url string
	MP3_Url     string
	Duration    int
	Description string
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func (pl Playlist) NotInStorage(stg storage.Storage) ([]string, error) {
	fs, err := stg.List()
	if err != nil {
		return nil, fmt.Errorf("listing files in storage: %w", err)
	}

	var notStored []string
	for _, e := range pl.Entries {
		if !contains(fs, e.Id+".mp3") {
			notStored = append(notStored, e.Id)
		}
	}
	return notStored, nil
}

func (pl Playlist) NotInPlaylist(stg storage.Storage) ([]string, error) {
	var plFs []string

	for _, e := range pl.Entries {
		plFs = append(plFs, e.Id+".mp3")
	}

	fs, err := stg.List()
	if err != nil {
		return nil, fmt.Errorf("listing files in storage: %w", err)
	}

	var notInPlaylist []string
	for _, f := range fs {
		if !contains(plFs, f) && f != "feed.rss" {
			notInPlaylist = append(notInPlaylist, f)
		}
	}

	return notInPlaylist, nil
}
