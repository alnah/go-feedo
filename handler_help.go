package main

import (
	"fmt"
)

// handlerHelp prints out the list of available commands and their usage
func handlerHelp(_ *state, _ command) error {
	fmt.Println("Available commands:")
	fmt.Println("=====================================")
	fmt.Println("register <name>           - Create a new user and set as current")
	fmt.Println("login <name>              - Set an existing user as the current user")
	fmt.Println("users                     - List all registered users")
	fmt.Println("reset                     - Delete all users from the database (dev only!)")
	fmt.Println("addfeed <name> <url>      - Add a new feed and follow it as current user")
	fmt.Println("feeds                     - List all available feeds")
	fmt.Println("follow <url>              - Follow an existing feed by URL")
	fmt.Println("unfollow <url>            - Unfollow a feed by URL")
	fmt.Println("following                 - List feeds followed by current user")
	fmt.Println("browse [limit]            - Browse posts from followed feeds (default limit is 2)")
	fmt.Println("agg [duration]            - Collect feeds once or every duration (e.g., 10s, 1m)")
	fmt.Println("help                      - Show this help message")
	fmt.Println("=====================================")
	return nil
}
