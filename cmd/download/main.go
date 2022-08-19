package main

import (
	"fmt"
	"os"

	"github.com/bitwombat/listen-later/video/youtube"
)

const numargs = 2

func main() {
	if len(os.Args) != numargs {
		fmt.Print(`Usage: download <YouTube videoId>
Example:
$ download v2349d34K
`)
		os.Exit(1)
	}

	videoId := os.Args[1]

	youtube.DownloadMP3(videoId)
}
