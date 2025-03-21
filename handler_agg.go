package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/alnah/go-gator/internal/database"
	"github.com/google/uuid"
)

// handlerAgg fetches the RSS feeds, parse them, and print the posts title in the console
// all in a long-running loop.
func handlerAgg(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		log.Printf("Collecting feeds...")
		if err := scrapeFeeds(s); err != nil {
			return err
		}
		return nil
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
// using its URL, create the posts for that feed into the database,
// and print the post titles in the console
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
		publishedAt := sql.NullTime{}
		if t, err := time.Parse(time.RFC1123Z, item.PubDate); err == nil {
			publishedAt = sql.NullTime{
				Time:  t,
				Valid: true,
			}
		}
		_, err = s.dbQr.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
			FeedID:      nextFeedToFetch.ID,
			Title:       item.Title,
			Description: item.Description,
			Url:         item.Link,
			PublishedAt: publishedAt,
		})
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key value") {
				continue
			}
			log.Printf("couldn't create post: %v", err)
			continue
		}
	}
	log.Printf("Feed %s collected, %v posts found", nextFeedToFetch.Name, len(feedData.Channel.Item))
	return nil
}
