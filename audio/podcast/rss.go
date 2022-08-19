package podcast

// Creates an RSS podcast feed

import (
	"bytes"
	"html"
	"text/template"
	"time"

	"github.com/bitwombat/listen-later/playlist"
)

const (
	rssFeedTemplate = `<?xml version="1.0" encoding="utf-8"?>
<rss xmlns:atom="http://www.w3.org/2005/Atom" xmlns:itunes="http://www.itunes.com/dtds/podcast-1.0.dtd" xmlns:itunesu="http://www.itunesu.com/feed" version="2.0">
<channel>
<link>https://www.bitwombat.com.au</link>
<language>en-us</language>
<copyright>&#xA9;2020</copyright>
<webMaster>greg@bitwombat.com.au (Greg Bell)</webMaster>
<managingEditor>greg@bitwombat.com.au (Greg Bell)</managingEditor>
<image>
   <url>https://bitwombat.com.au/pods/bw.png</url>
   <title>{{escape .Title}}</title>
   <link>https://bitwombat.com.au</link>
</image>
<itunes:owner>
   <itunes:name>Greg Bell</itunes:name>
   <itunes:email>greg@bitwombat.com.au</itunes:email>
</itunes:owner>
<itunes:category text="Education">
   <itunes:category text="Higher Education" />
</itunes:category>
<itunes:keywords>separate, by, comma, and, space</itunes:keywords>
<itunes:explicit>no</itunes:explicit>
<itunes:image href="http://bitwombat.com.au/pods/bw.png" />
<atom:link href="https://www.YourSite.com/feed.xml" rel="self" type="application/rss+xml" />
<pubDate>Fri, 05 Oct 2018 09:00:00 GMT</pubDate>
<title>{{escape .Title}}</title>
<itunes:author>College, school, or department owning the podcast</itunes:author>
<description>Nothing needs to be here</description>
<itunes:summary>No summary is necessary</itunes:summary>
<itunes:subtitle>Short description of the podcast - 255 character max.</itunes:subtitle>
<lastBuildDate>Fri, 05 Oct 2018 09:00:00 GMT</lastBuildDate>
{{range .Entries}}<item>
   <title>{{escape .Title}}</title>
   <description>{{escape .Description}}</description>
   <itunes:summary></itunes:summary>
   <itunes:subtitle></itunes:subtitle>
   <itunesu:category itunesu:code="112" />
   <enclosure url="{{.MP3_Url}}" type="audio/mpeg" length="{{.Duration}}" />
   <guid>{{.Webpage_Url}}</guid>
   <itunes:duration>{{.Duration}}</itunes:duration>
   <pubDate>{{pubDate .Upload_Date}}</pubDate>
</item>{{end}}
</channel>
</rss>`
)

func pubDate(ytFormat string) string {
	// Comes in as 20220513, leaves as Fri, 05 Oct 2018 09:00:00 GMT
	t, err := time.Parse("20060102", ytFormat)
	if err != nil {
		panic(err)
	}

	return t.Format("Mon, _2 Jan 2006 08:00:00 MST")
}

func RSSFromPlaylist(pl playlist.Playlist) (string, error) {
	funcMap := template.FuncMap{
		"pubDate": pubDate,
		"escape":  html.EscapeString,
	}

	tmpl, err := template.New("Podcast Feed").Funcs(funcMap).Parse(rssFeedTemplate)
	if err != nil {
		return "", err
	}

	var buf []byte
	feed := bytes.NewBuffer(buf)
	err = tmpl.Execute(feed, pl)

	if err != nil {
		return "", err
	}

	return feed.String() + "\n", nil // end file with newline because it might be expected
}
