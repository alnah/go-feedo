package main

import (
	"context"
	"fmt"
	"time"

	"github.com/alnah/go-feedo/internal/database"
	"github.com/google/uuid"
)

// handlerAddFeed adds a feed to the feeds table, and follows it for the current user who added the feed
func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 2 {
		return fmt.Errorf("usage: %v <name> <url>", cmd.name)
	}
	name := cmd.args[0]
	url := cmd.args[1]
	feed, err := s.dbQr.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		Name:      name,
		Url:       url,
	})
	if err != nil {
		return fmt.Errorf("couldn't create feed: %w", err)
	}
	feedFollow, err := s.dbQr.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("couldn't create feed follow: %w", err)
	}
	fmt.Println("Feed created successfully:")
	printFeed(feed, user)
	fmt.Println()
	fmt.Println("Feed followed successfully:")
	printFeedFollow(feedFollow.UserName, feedFollow.FeedName)
	fmt.Println("=====================================")
	return nil
}

// handlerListFeeds gets all the feeds from the feed table
func handlerListFeeds(s *state, cmd command) error {
	feeds, err := s.dbQr.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("couldn't get feeds: %w", err)
	}
	if len(feeds) == 0 {
		fmt.Println("No feeds found.")
		return nil
	}
	fmt.Printf("Found %d feeds:\n", len(feeds))
	for _, feed := range feeds {
		user, err := s.dbQr.GetUserById(context.Background(), feed.UserID)
		if err != nil {
			return fmt.Errorf("couldn't get user: %w", err)
		}
		printFeed(feed, user)
		fmt.Println("=====================================")
	}
	return nil
}

func printFeed(feed database.Feed, user database.User) {
	fmt.Printf("* ID:            %s\n", feed.ID)
	fmt.Printf("* Created:       %v\n", feed.CreatedAt)
	fmt.Printf("* Updated:       %v\n", feed.UpdatedAt)
	fmt.Printf("* Name:          %s\n", feed.Name)
	fmt.Printf("* URL:           %s\n", feed.Url)
	fmt.Printf("* User:          %s\n", user.Name)
	fmt.Printf("* LastFetchedAt: %v\n", feed.LastFetchedAt.Time)
}
