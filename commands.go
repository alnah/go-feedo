package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/alnah/go-gator/internal/config"
	"github.com/alnah/go-gator/internal/database"
	"github.com/google/uuid"
)

// state gives to the handlers an access to the application state and the database queries
type state struct {
	dbQr  *database.Queries
	dbCfg *config.DatabaseConfig
}

// command contains a command name and a slice of arguments expected by the handlers
type command struct {
	name string
	args []string
}

// commandHandler is the function signature of all command handlers
type commandHandler func(s *state, cmd command) error

// commands holds all the commands the CLI can handle in a registry
type commands struct {
	registry map[string]commandHandler
}

// register adds a new command handler to the registry if it doesn't exist
func (c *commands) register(name string, f commandHandler) {
	if c.registry == nil {
		c.registry = make(map[string]commandHandler)
	}
	if _, exist := c.registry[name]; !exist {
		c.registry[name] = f
	}
}

// run retrieves a given command and run it with the provided state if it exists
func (c *commands) run(s *state, cmd command) error {
	if handler, exist := c.registry[cmd.name]; exist {
		err := handler(s, cmd)
		if err != nil {
			return err
		}
	}
	return nil
}

// handleLogin sets the current user name of the database configuration to the given username
func handleLogin(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return errors.New("Please provide a single username, e.g., \"login alice\"")
	}
	username := cmd.args[0]
	_, err := s.dbQr.GetUserByName(context.Background(), username)
	if err != nil {
		return fmt.Errorf("Error getting user from the database %w:", err)
	}
	if err := s.dbCfg.SetUser(username); err != nil {
		return err
	}
	fmt.Printf("The user has been set to: %q", username)
	return nil
}

// handleAddUser creates a new username and insert it into the users table
// it also sets the new registered user as the current user name of the database configuration
func handleAddUser(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return errors.New("Please provide a single username, e.g.,\"register alice\"")
	}
	username := cmd.args[0]
	_, err := s.dbQr.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: sql.NullTime{Time: time.Now(), Valid: true},
		UpdatedAt: sql.NullTime{Time: time.Now(), Valid: true},
		Name:      username,
	})
	if err != nil {
		return fmt.Errorf("Error inserting user into the database: %q", err)
	}
	if err = s.dbCfg.SetUser(username); err != nil {
		return err
	}
	fmt.Printf("The user has been set to: %q", username)
	return nil
}

// handleResetAllUsers delete all the users from the users table
// it's really to work on CRUD operations, of course never do that into a real app...
func handleResetAllUsers(s *state, cmd command) error {
	if err := s.dbQr.DeleteAllUsers(context.Background()); err != nil {
		return fmt.Errorf("Error deleting all users from the databse: %q", err)
	}
	return nil
}

// handleListAllUsers retrieve all the first 100 users from users table
func handleListAllUsers(s *state, _ command) error {
	users, err := s.dbQr.GetAllUsers(context.Background(), database.GetAllUsersParams{Limit: 100, Offset: 0})
	if err != nil {
		return fmt.Errorf("Error retrieving all the users from the database %q", err)
	}
	var output strings.Builder
	for _, user := range users {
		if user.Name == s.dbCfg.CurrentUserName {
			output.WriteString(fmt.Sprintf("* %s (current)\n", user.Name))
			continue
		}
		output.WriteString(fmt.Sprintf("* %s\n", user.Name))
	}
	fmt.Print(output.String())
	return nil
}

// handleAggregate aggregates feeds from one RSS feed URL
func handleAggregate(_ *state, _ command) error {
	// if len(cmd.args) != 1 {
	// 	return errors.New("Please provide one RSS feed URL, e.g.\"register https://www.myfeed.com/index.xml\"")
	// }
	feedURL := "https://www.wagslane.dev/index.xml" // cmd.args[0]
	rssFeed, err := fetchFeed(context.Background(), feedURL)
	if err != nil {
		return fmt.Errorf("Error aggregating feeds: %w", err)
	}
	fmt.Printf("%+v", rssFeed)
	return nil
}

// handleAddFeed adds a feed to the feeds table
func handleAddFeed(s *state, cmd command) error {
	if len(cmd.args) != 2 {
		return errors.New("Please provide a feed name, and its url, e.g., \"addfeed 'Hacker News RSS' 'https://hnrss.org/newest'\"")
	}
	ctx, feedName, feedUrl := context.Background(), cmd.args[0], cmd.args[1]
	user, err := s.dbQr.GetUserByName(ctx, s.dbCfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("Error getting user from the database: %q", err)
	}
	feed, err := s.dbQr.CreateFeed(ctx, database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: sql.NullTime{Time: time.Now(), Valid: true},
		UpdatedAt: sql.NullTime{Time: time.Now(), Valid: true},
		Name:      feedName,
		Url:       feedUrl,
		UserID:    user.ID,
	})
	if err != nil {
		return fmt.Errorf("Error inserting feed into the database: %q", err)
	}
	fmt.Printf("%+v", feed)
	return nil
}

// handleListAllFeeds gets all the feeds from the feed table
func handleListAllFeeds(s *state, _ command) error {
	ctx := context.Background()
	feeds, err := s.dbQr.GetAllFeeds(ctx, database.GetAllFeedsParams{Limit: 100, Offset: 0})
	if err != nil {
		return fmt.Errorf("Error getting all feeds from the database: %q", err)
	}
	for _, feed := range feeds {
		usern, err := s.dbQr.GetUserByID(ctx, feed.UserID)
		if err != nil {
			return fmt.Errorf("Error getting user by ID from the database: %w", err)
		}
		fmt.Printf("* %s: %s (owner: %s)\n", feed.Name, feed.Url, usern.Name)
	}
	return nil
}
