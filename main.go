package main

import (
	"log"
	"os"

	"github.com/lordvorath/gator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}
	userState := state{cfg: &cfg}
	availableCommands := commands{commandMap: make(map[string]func(*state, command) error)}
	availableCommands.register("login", handlerLogin)
	args := os.Args
	if len(args) < 2 {
		log.Fatal("Not enough arguments")
	}
	cmd := command{name: args[1], args: args[2:]}
	err = availableCommands.run(&userState, cmd)
	if err != nil {
		log.Fatalf("error encountered: %v", err)
	}
}
