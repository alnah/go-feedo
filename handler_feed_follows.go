package main

import (
	"context"
	"fmt"
	"time"

	"github.com/alnah/go-gator/internal/database"
	"github.com/google/uuid"
)

// handlerFollow creates a new feed follow record for the current user in the feed_follows table
func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("usage: %v <feed_url>", cmd.name)
	}
	feed, err := s.dbQr.GetFeedByURL(context.Background(), cmd.args[0])
	if err != nil {
		return fmt.Errorf("couldn't get feed: %w", err)
	}
	ffRow, err := s.dbQr.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("couldn't create feed follow: %w", err)
	}
	fmt.Println("Feed follow created:")
	printFeedFollow(ffRow.UserName, ffRow.FeedName)
	return nil
}

// handlerListFeedFollows retrieves all the names of the feeds the current user is following
func handlerListFeedFollows(s *state, cmd command, user database.User) error {
	feedFollows, err := s.dbQr.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("couldn't get feed follows: %w", err)
	}
	if len(feedFollows) == 0 {
		fmt.Println("No feed follows found for this user.")
		return nil
	}
	fmt.Printf("Feed follows for user %s:\n", user.Name)
	for _, ff := range feedFollows {
		fmt.Printf("* %s\n", ff.FeedName)
	}
	return nil
}

func printFeedFollow(username, feedname string) {
	fmt.Printf("* User:          %s\n", username)
	fmt.Printf("* Feed:          %s\n", feedname)
}

// handlerUnfollow allows to unfollow a feed for the current user
func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("usage: %s <url>", cmd.name)
	}
	url := cmd.args[0]
	err := s.dbQr.DeleteFeedFollow(context.Background(), database.DeleteFeedFollowParams{UserID: user.ID, Url: url})
	if err != nil {
		return fmt.Errorf("couldn't delete feed follow: %w", err)
	}
	return nil
}
