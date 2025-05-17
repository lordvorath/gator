package main

import (
	"log"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	cfg, dbQueries := middlewareInit()
	userState := state{cfg: cfg, db: dbQueries}

	availableCommands := commands{commandMap: make(map[string]func(*state, command) error)}
	availableCommands.register("login", handlerLogin)
	availableCommands.register("register", handlerRegister)
	availableCommands.register("reset", handlerReset)
	availableCommands.register("users", handlerUSers)
	availableCommands.register("agg", handlerAggregator)
	availableCommands.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	availableCommands.register("feeds", handlerListFeeds)
	availableCommands.register("follow", middlewareLoggedIn(handlerAddFollow))
	availableCommands.register("following", middlewareLoggedIn(handlerListFollows))
	availableCommands.register("unfollow", middlewareLoggedIn(handlerDeleteFollow))

	args := os.Args
	if len(args) < 2 {
		log.Fatal("not enough arguments")
	}
	cmd := command{name: args[1], args: args[2:]}
	err := availableCommands.run(&userState, cmd)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
}
