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

// RSSFeed represents an RSS feed parsed from XML
type RSSFeed struct {
	// Channel contains metadata and items of the RSS feed
	Channel struct {
		// Title is the feed's title
		Title string `xml:"title"`
		// Link is the URL to the feed
		Link string `xml:"link"`
		// Description provides details about the feed
		Description string `xml:"description"`
		// Item is the list of feed entries
		Item []RSSItem `xml:"item"`
	} `xml:"channel"`
}

// RSSItem represents an individual entry in an RSS feed
type RSSItem struct {
	// Title of the feed item
	Title string `xml:"title"`
	// Link to the full content of the item
	Link string `xml:"link"`
	// Description or summary of the item
	Description string `xml:"description"`
	// PubDate is the publication date of the item
	PubDate string `xml:"pubDate"`
}

// fetchFeed retrieves and parses an RSS feed from the specified URL using the provided context.
func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return nil, fmt.Errorf("Error creating a new request: %q", err)
	}
	req.Header.Set("User-Agent", "go-feedo")
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
