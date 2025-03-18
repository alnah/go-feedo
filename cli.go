package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/alnah/go-gator/internal/config"
	"github.com/alnah/go-gator/internal/database"
)

// initClit initialize the command-line client and should be used from main entry point
// it orchestrates the different parts of the program together
func initCli() {
	dbCfg, err := config.Read()
	if err != nil {
		log.Fatalf("Error reading config: %v", err)
	}
	dbCon, err := sql.Open("postgres", dbCfg.URL)
	if err != nil {
		log.Fatalf("Error opening PostgreSQL database: %v", err)
	}
	dbQr := database.New(dbCon)
	if err != nil {
		log.Fatalf("Error")
	}
	s := &state{dbCfg: &dbCfg, dbQr: dbQr}
	cmds := commands{}
	cmds.register("login", handleLogin)
	cmds.register("register", handleRegister)
	cmds.register("reset", handleReset)
	if len(os.Args) < 2 {
		log.Fatal("Usage: cli <command> [args...]")
	}
	cmdName, cmdArgs := os.Args[1], os.Args[2:]
	cmd := command{name: cmdName, args: cmdArgs}
	err = cmds.run(s, cmd)
	if err != nil {
		log.Fatal(err)
	}
}
