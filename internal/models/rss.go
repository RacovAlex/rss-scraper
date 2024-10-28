package models

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"time"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Language    string    `xml:"language"`
		Item        []RSSItem `xml:"item,omitempty"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func UrlToFeed(url string) (RSSFeed, error) {
	const fn = "models.UrlToFeed"
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	resp, err := client.Get(url)
	if err != nil {
		return RSSFeed{}, fmt.Errorf("%s: %w", fn, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return RSSFeed{}, fmt.Errorf("%s: %w", fn, err)
	}
	rssFeed := RSSFeed{}
	if err := xml.Unmarshal(body, &rssFeed); err != nil {
		return RSSFeed{}, fmt.Errorf("%s: %w", fn, err)
	}
	return rssFeed, nil
}
