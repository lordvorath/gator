package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/lordvorath/gator/internal/config"
	"github.com/lordvorath/gator/internal/database"

	_ "github.com/lib/pq"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}
	db, err := sql.Open("postgres", cfg.DbUrl)
	if err != nil {
		log.Fatal("could not open database")
	}
	dbQueries := database.New(db)
	userState := state{cfg: &cfg, db: dbQueries}

	availableCommands := commands{commandMap: make(map[string]func(*state, command) error)}
	availableCommands.register("login", handlerLogin)
	availableCommands.register("register", handlerRegister)
	availableCommands.register("reset", handlerReset)
	availableCommands.register("users", handlerUSers)
	availableCommands.register("agg", handlerAggregator)
	availableCommands.register("addfeed", handlerAddFeed)

	args := os.Args
	if len(args) < 2 {
		log.Fatal("not enough arguments")
	}
	cmd := command{name: args[1], args: args[2:]}
	err = availableCommands.run(&userState, cmd)
	if err != nil {
		log.Fatalf("error encountered: %v", err)
	}
}
