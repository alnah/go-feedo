package main

import (
	"context"
	"fmt"
	"log"
)

// handlerReset delete all the users from the users table
// it's really to work on CRUD operations, of course never do that into a real app...
func handlerReset(s *state, cmd command) error {
	err := s.dbQr.DeleteUsers(context.Background())
	if err != nil {
		return fmt.Errorf("couldn't delete users: %w", err)
	}
	log.Println("Database reset successfully!")
	return nil
}
