package main

import (
	"context"
	"fmt"
	"time"

	"github.com/alnah/go-feedo/internal/database"
	"github.com/google/uuid"
)

// handlerRegister creates a new username and insert it into the users table
// it also sets the new registered user as the current user name of the database configuration
func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("usage: %v <name>", cmd.name)
	}
	name := cmd.args[0]
	user, err := s.dbQr.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      name,
	})
	if err != nil {
		return fmt.Errorf("couldn't create user: %w", err)
	}
	err = s.dbCfg.SetUser(user.Name)
	if err != nil {
		return fmt.Errorf("couldn't set current user: %w", err)
	}
	fmt.Println("User created successfully:")
	printUser(user)
	return nil
}

// handlerLogin sets the current user name of the database configuration to the given username
func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("usage: %v <name>", cmd.name)
	}
	name := cmd.args[0]
	_, err := s.dbQr.GetUser(context.Background(), name)
	if err != nil {
		return fmt.Errorf("couldn't find user: %w", err)
	}
	err = s.dbCfg.SetUser(name)
	if err != nil {
		return fmt.Errorf("couldn't set current user: %w", err)
	}
	fmt.Println("User switched successfully!")
	return nil
}

// handlerListUsers retrieves all the users from the users table
func handlerListUsers(s *state, cmd command) error {
	users, err := s.dbQr.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("couldn't list users: %w", err)
	}
	for _, user := range users {
		if user.Name == s.dbCfg.CurrentUserName {
			fmt.Printf("* %v (current)\n", user.Name)
			continue
		}
		fmt.Printf("* %v\n", user.Name)
	}
	return nil
}

func printUser(user database.User) {
	fmt.Printf(" * ID:      %v\n", user.ID)
	fmt.Printf(" * Name:    %v\n", user.Name)
}
