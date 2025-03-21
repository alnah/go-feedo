package main

import (
	"context"
	"fmt"

	"github.com/alnah/go-feedo/internal/database"
)

// middlewareLoggedIn allows to change the function signature of our handlers that
// require a logged in user to accept a user as an argument and DRY up the code
func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) commandHandler {
	return func(s *state, cmd command) error {
		user, err := s.dbQr.GetUser(context.Background(), s.dbCfg.CurrentUserName)
		if err != nil {
			return fmt.Errorf("couldn't get the current user: %w", err)
		}
		return handler(s, cmd, user)
	}
}
