package main

import (
	"fmt"

	"github.com/lordvorath/gator/internal/config"
)

type state struct {
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
		return fmt.Errorf("%v: inccorrect number of arguments: expected 1, got %d", cmd.name, len(cmd.args))
	}
	err := s.cfg.SetUser(cmd.args[0])
	if err != nil {
		return err
	}
	fmt.Printf("User set to: %v\n", cmd.args[0])
	return nil
}
