package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/alnah/go-gator/internal/database"
)

// handlerAgg fetches the RSS feeds, parse them, and print the posts title in the console
// all in a long-running loop.
func handlerAgg(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("usage: %v <time_between_reqs>", cmd.name)
	}
	timeBtwReqs, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return fmt.Errorf("invalid duration: %w", err)
	}
	log.Printf("Collecting feeds every %s...", timeBtwReqs.String())
	ticker := time.NewTicker(timeBtwReqs)
	for ; ; <-ticker.C {
		if err := scrapeFeeds(s); err != nil {
			return err
		}
	}
}

// scrapeFeeds goes to the next feed to fetch, mark it as fetched, fetch the feed data
// using its URL, and print the post titles in the console
func scrapeFeeds(s *state) error {
	ctx := context.Background()
	nextFeedToFetch, err := s.dbQr.GetNextFeedToFetch(ctx)
	if err != nil {
		return fmt.Errorf("couldn't get next feed to fetch: %w", err)
	}
	log.Println("Found a feed to fetch!")
	err = s.dbQr.MarkFeedFetched(ctx, database.MarkFeedFetchedParams{
		ID:        nextFeedToFetch.ID,
		UpdatedAt: time.Now(),
	})
	if err != nil {
		return fmt.Errorf("couldn't mark feed as fetched: %w", err)
	}
	feedData, err := fetchFeed(ctx, nextFeedToFetch.Url)
	if err != nil {
		return fmt.Errorf("couldn't fetch feed: %w", err)
	}
	for _, item := range feedData.Channel.Item {
		fmt.Printf("Feed: %+v\n", item.Title)
	}
	return err
}
