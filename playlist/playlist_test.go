package playlist

import (
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

type mockStorage struct{}

var files []string //nolint

func (m mockStorage) List() ([]string, error) {
	return files, nil
}

func (m mockStorage) FileExists(string) (bool, error)        { return true, nil }
func (m mockStorage) Get(string) (string, error)             { return "", nil }
func (m mockStorage) Put(io.Reader, string) error            { return nil }
func (m mockStorage) PutIfNotExists(io.Reader, string) error { return nil }
func (m mockStorage) Delete(name string) error               { return nil }
func (m mockStorage) DeleteList(l []string) error            { return nil }

func Test_DeletedFromYouTube(t *testing.T) {
	tt := []struct {
		purpose  string
		files    []string
		entries  []string
		toDelete []string
	}{
		{
			purpose:  "Some files are not in entries",
			files:    []string{"aaa.mp3", "bbb.mp3", "ddd.mp3"},
			entries:  []string{"bbb"},
			toDelete: []string{"aaa.mp3", "ddd.mp3"},
		},
		{
			purpose:  "Some entries are not in files",
			files:    []string{"aaa.mp3", "bbb.mp3", "ddd.mp3"},
			entries:  []string{"aaa", "bbb", "ccc", "ddd", "eee"},
			toDelete: []string{},
		},
		{
			purpose:  "All files are not in entries",
			files:    []string{"ddd.mp3", "eee.mp3", "fff.mp3"},
			entries:  []string{"aaa", "bbb", "ccc"},
			toDelete: []string{"ddd.mp3", "eee.mp3", "fff.mp3"},
		},
		{
			purpose:  "All files are in entries",
			files:    []string{"aaa.mp3", "bbb.mp3", "ccc.mp3"},
			entries:  []string{"aaa", "bbb", "ccc"},
			toDelete: []string{},
		},
		{
			purpose:  "Files are empty",
			files:    []string{},
			entries:  []string{"aaa", "bbb", "ccc"},
			toDelete: []string{},
		},
		{
			purpose:  "Entries are empty",
			files:    []string{"aaa.mp3", "bbb.mp3", "ccc.mp3"},
			entries:  []string{},
			toDelete: []string{"aaa.mp3", "bbb.mp3", "ccc.mp3"},
		},
		{
			purpose:  "feed.rss is left alone",
			files:    []string{"aaa.mp3", "bbb.mp3", "ccc.mp3", "feed.rss"},
			entries:  []string{},
			toDelete: []string{"aaa.mp3", "bbb.mp3", "ccc.mp3"},
		},
		{
			purpose:  "Other types of files are mopped up",
			files:    []string{"aaa.mp3", "bbb.mp3", "ccc.mp3", "feed.barf"},
			entries:  []string{},
			toDelete: []string{"aaa.mp3", "bbb.mp3", "ccc.mp3", "feed.barf"},
		},
	}

	for _, testcase := range tt {
		t.Run("different sets", func(t *testing.T) {
			files = testcase.files
			pl := &Playlist{}

			for _, e := range testcase.entries {
				pl.Entries = append(pl.Entries, Entry{Id: e})
			}

			es, err := pl.NotInPlaylist(mockStorage{})
			if err != nil {
				t.Fatal("Should not have returned nil")
			}

			require.Len(t, es, len(testcase.toDelete), "Entries in storage have been deleted from YouTube")
			for i, r := range testcase.toDelete {
				require.Equal(t, r, es[i], "Entries deleted")
			}
		})
	}
}

func Test_DownloadFromYouTube(t *testing.T) {
	testTable := []struct {
		purpose    string
		files      []string
		entries    []string
		toDownload []string
	}{
		{
			purpose:    "Some files are not in entries",
			files:      []string{"aaa.mp3", "bbb.mp3", "ddd.mp3"},
			entries:    []string{"bbb"},
			toDownload: []string{},
		},
		{
			purpose:    "Some entries are not in files",
			files:      []string{"aaa.mp3", "bbb.mp3", "ddd.mp3"},
			entries:    []string{"aaa", "bbb", "ccc", "ddd", "eee"},
			toDownload: []string{"ccc", "eee"},
		},
		{
			purpose:    "All files are not in entries",
			files:      []string{"ddd.mp3", "eee.mp3", "fff.mp3"},
			entries:    []string{"aaa", "bbb", "ccc"},
			toDownload: []string{"aaa", "bbb", "ccc"},
		},
		{
			purpose:    "All files are in entries",
			files:      []string{"aaa.mp3", "bbb.mp3", "ccc.mp3"},
			entries:    []string{"aaa", "bbb", "ccc"},
			toDownload: []string{},
		},
		{
			purpose:    "Files are empty",
			files:      []string{},
			entries:    []string{"aaa", "bbb", "ccc"},
			toDownload: []string{"aaa", "bbb", "ccc"},
		},
		{
			purpose:    "Entries are empty",
			files:      []string{"aaa.mp3", "bbb.mp3", "ccc.mp3"},
			entries:    []string{},
			toDownload: []string{},
		},
	}

	for _, testcase := range testTable {
		t.Run("different sets", func(t *testing.T) {
			files = testcase.files
			pl := &Playlist{}

			for _, e := range testcase.entries {
				pl.Entries = append(pl.Entries, Entry{Id: e})
			}

			es, err := pl.NotInStorage(mockStorage{})
			if err != nil {
				t.Fatal("Should not have returned nil")
			}

			require.Len(t, es, len(testcase.toDownload), "Entries not in storage need to be downloaded")
			for i, r := range testcase.toDownload {
				require.Equal(t, r, es[i], "Entries downloaded")
			}
		})
	}
}
