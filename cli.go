package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lordvorath/gator/internal/config"
	"github.com/lordvorath/gator/internal/database"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}

type command struct {
	name string
	args []string
}

type commands struct {
	commandMap map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	f, ok := c.commandMap[cmd.name]
	if !ok {
		availableCommands := make([]string, len(c.commandMap))
		i := 0
		for k := range c.commandMap {
			availableCommands[i] = k
			i++
		}
		return fmt.Errorf("unrecognized command, available commands: %v", availableCommands)
	}
	err := f(s, cmd)
	if err != nil {
		return err
	}
	return nil
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.commandMap[name] = f
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("usage: %v <name>", cmd.name)
	}
	usr, err := s.db.GetUser(context.Background(), cmd.args[0])
	if err != nil {
		return fmt.Errorf("user %v does not exist: %w", cmd.args[0], err)
	}
	err = s.cfg.SetUser(usr.Name)
	if err != nil {
		return err
	}
	fmt.Printf("User set to: %v\n", cmd.args[0])
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("usage: %v <name>", cmd.name)
	}
	newUser := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.args[0],
	}
	usr, err := s.db.CreateUser(context.Background(), newUser)
	if err != nil {
		return fmt.Errorf("error during user creation: %w", err)
	}
	err = s.cfg.SetUser(usr.Name)
	if err != nil {
		return fmt.Errorf("error updating config: %w", err)
	}
	fmt.Printf("User created successfully:\nUUID: %v\nCreated: %v\nUpdated: %v\nName: %v\n",
		usr.ID, usr.CreatedAt, usr.UpdatedAt, usr.Name)
	return nil
}
