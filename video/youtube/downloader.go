package youtube

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"

	"github.com/bitwombat/listen-later/storage"
)

// DownloadMP3ToStorage downloads a video by ID, and saves the resulting MP3 to storage.
func DownloadMP3ToStorage(videoID string, stg storage.Storage) error {
	r, err := DownloadMP3(videoID)
	if err != nil {
		return err
	}

	if err := stg.PutIfNotExists(r, videoID+".mp3"); err != nil {
		return err
	}

	// Remove local one
	err = os.Remove(videoID + ".mp3")
	if err != nil {
		log.Fatalf("failed deleting local mp3 file: %v", err)
	}

	return nil
}

// DownloadMP3 downloads a youtube video as an MP3 to disk, and provides a reader to that file.
func DownloadMP3(videoID string) (io.Reader, error) {
	URL := "https://www.youtube.com/watch?v=" + videoID

	log.Println("Downloading " + URL)

	cmd1 := exec.Command("yt-dlp",
		"--extract-audio",
		"--embed-thumbnail",
		"--output", "%(id)s.%(ext)s",
		"--no-progress",
		"--no-colors",
		"--ignore-errors",
		"--audio-format", "mp3",
		URL)
	cmd1.Stderr = os.Stderr

	if err := cmd1.Run(); err != nil {
		return nil, fmt.Errorf("failed downloading mp3: %v", err)
	}

	log.Println("Completed download")

	r, err := os.Open(videoID + ".mp3")
	if err != nil {
		return nil, fmt.Errorf("failed opening mp3 file: %v", err)
	}

	return r, nil
}
