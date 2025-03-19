package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"time"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return nil, fmt.Errorf("Error creating a new request: %q", err)
	}
	req.Header.Set("User-Agent", "gator")
	req.Header.Set("Content-Type", "application/rss+xml")
	client := &http.Client{Timeout: 10 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Error getting response: %q", err)
	}
	defer func() { _ = res.Body.Close() }()
	byt, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("Error reading response body: %q", err)
	}
	var rssFeed RSSFeed
	if err = xml.Unmarshal(byt, &rssFeed); err != nil {
		return nil, fmt.Errorf("Error unmarshalling: %q", err)
	}
	var toUnescape []string
	toUnescape = append(toUnescape, rssFeed.Channel.Title, rssFeed.Channel.Description)
	for _, item := range rssFeed.Channel.Item {
		toUnescape = append(toUnescape, item.Description, item.Title)
	}
	for i := range toUnescape {
		toUnescape[i] = html.UnescapeString(toUnescape[i])
	}
	return &rssFeed, nil
}
