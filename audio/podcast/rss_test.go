package podcast

/* Test the RSS feed generator
 */

import (
	"os"
	"testing"

	"github.com/bitwombat/listen-later/playlist"
	"github.com/stretchr/testify/require"
)

func Test_template(t *testing.T) {
	t.Parallel()
	t.Run("Should format publish date correctly", func(t *testing.T) {
		t.Parallel()
		// GIVEN: An upload date as read from the YouTube playlist
		uploadDate := "20220513"

		// WHEN: We format that as a date for the RSS feed
		got := pubDate(uploadDate)

		// THEN: It should be in the format expected by a podcast feed
		require.Equal(t, "Fri, 13 May 2022 08:00:00 UTC", got, "RSS feed time format")
	})

	t.Run("Should fill in template from list of podcasts", func(t *testing.T) {
		t.Parallel()
		// GIVEN: A slice of Podcasts and a template
		entries := []playlist.Entry{
			{
				Id:          "AAAAAAAA",
				Title:       "My Title",
				Upload_Date: "20220525",
				MP3_Url:     "https://www.youtube.com/watch?v=AAAAAAAAAAA",
				Webpage_Url: "https://www.youtube.com/watch?v=AAAAAAAAAAA",
				Duration:    121,
				Description: "some description here",
			}, {
				Id:          "ZZZZZZZZ",
				Title:       "Ooo, one more",
				Upload_Date: "20220526",
				MP3_Url:     "https://www.youtube.com/watch?v=ZZZZZZZZZZZ",
				Webpage_Url: "https://www.youtube.com/watch?v=ZZZZZZZZZZZ",
				Duration:    60,
				Description: "riveting episode",
			},
		}
		list := playlist.Playlist{
			Title:   "Greg's Listen Later List",
			Entries: entries,
		}

		// WHEN: We generate the RSS from it
		rss, err := RSSFromPlaylist(list)
		if err != nil {
			t.Fail()
		}

		// THEN: The output should match our golden file
		golden, err := os.ReadFile("feed_golden.rss")
		if err != nil {
			panic(err)
		}
		require.Equal(t, string(golden), string(rss), "rss output")
	})
}
