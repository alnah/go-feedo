package main

import (
	"log"
	"os"

	"github.com/alnah/go-gator/internal/config"
)

// initClit initialize the command-line client and should be used from main entry point
// it orchestrates the different parts of the program together
func initCli() {
	db, err := config.Read()
	if err != nil {
		log.Fatalf("Error reading config: %v", err)
	}
	s := &state{&db}
	cmds := commands{}
	cmds.register("login", handleLogin)
	if len(os.Args) <= 2 {
		log.Fatal("Usage: cli <command> [args...]")
	}
	cmdName, cmdArgs := os.Args[1], os.Args[2:]
	cmd := command{name: cmdName, args: cmdArgs}
	err = cmds.run(s, cmd)
	if err != nil {
		log.Fatal(err)
	}
}
