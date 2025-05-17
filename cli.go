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
	usr, err := s.db.GetUserByName(context.Background(), cmd.args[0])
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

func handlerReset(s *state, cmd command) error {
	if len(cmd.args) != 0 {
		return fmt.Errorf("usage: %v", cmd.name)
	}
	err := s.db.Reset(context.Background())
	if err != nil {
		return fmt.Errorf("error during users DB reset: %w", err)
	}
	err = s.db.ResetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("error during feeds DB reset: %w", err)
	}
	err = s.db.ResetFeedFollows(context.Background())
	if err != nil {
		return fmt.Errorf("error during feeds DB reset: %w", err)
	}
	fmt.Println("database reset successfully")
	return nil
}

func handlerUSers(s *state, cmd command) error {
	if len(cmd.args) != 0 {
		return fmt.Errorf("usage: %v", cmd.name)
	}
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("could not retriever users: %w", err)
	}
	fmt.Println("Users:")
	for _, u := range users {
		log := "* " + u.Name
		if s.cfg.CurrentUserName == u.Name {
			log = log + " (current)"
		}
		fmt.Println(log)
	}
	return nil
}

func handlerAggregator(s *state, cmd command) error {
	if len(cmd.args) != 0 {
		return fmt.Errorf("usage: %v", cmd.name)
	}
	feedUrl := "https://www.wagslane.dev/index.xml"
	rssFeed, err := fetchFeed(context.Background(), feedUrl)
	if err != nil {
		return fmt.Errorf("fetchFeed failed: %w", err)
	}
	fmt.Println(rssFeed.Channel.Title)
	fmt.Println(rssFeed.Channel.Description)
	fmt.Println(rssFeed.Channel.Link)
	fmt.Printf("%v\n", rssFeed.Channel.Item)
	return nil
}

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 2 {
		return fmt.Errorf("usage: %v <name> <url>", cmd.name)
	}
	newFeed := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.args[0],
		Url:       cmd.args[1],
		UserID:    user.ID,
	}
	f, err := s.db.CreateFeed(context.Background(), newFeed)
	if err != nil {
		return fmt.Errorf("failed to create feed: %w", err)
	}
	fmt.Printf("RSS Feed: %v added to database\n", f.Name)
	err = handlerAddFollow(s, command{name: "follow", args: []string{cmd.args[1]}}, user)
	if err != nil {
		return fmt.Errorf("failed to follow newly added feed: %w", err)
	}
	return nil
}

func handlerListFeeds(s *state, cmd command) error {
	if len(cmd.args) != 0 {
		return fmt.Errorf("usage: %v", cmd.name)
	}
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("failed to retrieve feeds: %w", err)
	}
	for _, feed := range feeds {
		u, err := s.db.GetUserById(context.Background(), feed.UserID)
		if err != nil {
			return fmt.Errorf("failed to retrieve user: %w", err)
		}
		fmt.Printf("Name: %v\n", feed.Name)
		fmt.Printf("URL: %v\n", feed.Url)
		fmt.Printf("Added by: %v\n", u.Name)
	}
	return nil
}

func handlerAddFollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("usage: %v <url>", cmd.name)
	}
	feed, err := s.db.GetFeedByURL(context.Background(), cmd.args[0])
	if err != nil {
		return fmt.Errorf("failed to find feed: %w", err)
	}
	newFollow := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}
	follow, err := s.db.CreateFeedFollow(context.Background(), newFollow)
	if err != nil {
		return fmt.Errorf("failed to create feed follow: %w", err)
	}
	fmt.Printf("User: %v is now following Feed: %v\n", follow.Username, follow.Feedname)
	return nil
}

func handlerListFollows(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 0 {
		return fmt.Errorf("usage: %v", cmd.name)
	}
	follows, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("failed to get follows for current user: %w", err)
	}
	fmt.Println("List of currently followed feeds:")
	for _, follow := range follows {
		fmt.Printf("%v\n", follow.Feedname)
	}
	return nil
}

func handlerDeleteFollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("usage: %v <url>", cmd.name)
	}
	toDelete := database.DeleteFeedFollowParams{
		Url:  cmd.args[0],
		Name: s.cfg.CurrentUserName,
	}
	err := s.db.DeleteFeedFollow(context.Background(), toDelete)
	if err != nil {
		return fmt.Errorf("could not delete feed: %w", err)
	}
	fmt.Println("Feed deleted successfully")
	return nil
}
