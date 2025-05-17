package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/lordvorath/gator/internal/config"
	"github.com/lordvorath/gator/internal/database"
)

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		currentUser, err := s.db.GetUserByName(context.Background(), s.cfg.CurrentUserName)
		if err != nil {
			return fmt.Errorf("failed to log in user: %w", err)
		}
		return handler(s, cmd, currentUser)
	}
}

func middlewareInit() (*config.Config, *database.Queries) {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}
	db, err := sql.Open("postgres", cfg.DbUrl)
	if err != nil {
		log.Fatal("could not open database")
	}
	dbQueries := database.New(db)
	return &cfg, dbQueries
}
