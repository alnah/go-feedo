package main

import (
	"github.com/alnah/go-feedo/internal/config"
	"github.com/alnah/go-feedo/internal/database"
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
