package main

import (
	"github.com/mmcdole/gofeed"
)

const feedURL = "https://news.google.com/rss?hl=pt&gl=BR&ceid=BR:pt"

// FeedReader ..
type FeedReader struct {
	fp    *gofeed.Parser
	feeds *gofeed.Feed
	i     int
}

// NewReader ..
func NewReader() *FeedReader {
	fp := gofeed.NewParser()
	feeds, _ := fp.ParseURL(feedURL)

	return &FeedReader{fp: fp, feeds: feeds, i: 0}
}

// GetFeed ...
func (fr *FeedReader) GetFeed() string {
	if fr.i >= fr.feeds.Len() {
		fr.feeds, _ = fr.fp.ParseURL(feedURL)
		fr.i = 0
	}

	title := fr.feeds.Items[fr.i].Title

	fr.i++

	return title

}

// GetURL ..
func (fr *FeedReader) GetURL() string {
	i := 0

	if fr.i > 0 {
		i = fr.i - 1
	}

	return fr.feeds.Items[i].Link
}

func (fr *FeedReader) PrevItem() string {
	i := 0

	if fr.i > 1 {
		i = fr.i - 2
	}

	fr.i = i

	return ""
}
