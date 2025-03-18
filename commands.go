package main

import (
	"errors"
	"fmt"
	"github.com/alnah/go-gator/internal/config"
)

// state gives to the handlers an access to the application state
type state struct{ *config.DatabaseConfig }

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
		return errors.New("Please provide one username, e.g.\"login alice\"")
	}
	username := cmd.args[0]
	fmt.Println(username)
	if err := s.SetUser(username); err != nil {
		return err
	}
	fmt.Printf("The user has been set to %q\n", username)
	return nil
}
