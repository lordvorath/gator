package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
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

func scrapeFeeds(s *state) error {
	//fmt.Println("finding next feed to fetch")
	nextFeed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return fmt.Errorf("could not find next feed to fetch: %w", err)
	}
	//fmt.Printf("fetching feed %v\n", nextFeed.Url)
	feed, err := fetchFeed(context.Background(), nextFeed.Url)
	if err != nil {
		return fmt.Errorf("failed to fetch feed: %w", err)
	}
	//fmt.Printf("marking feed fetched: %v", nextFeed.Name)
	err = s.db.MarkFeedFetched(context.Background(),
		database.MarkFeedFetchedParams{ID: nextFeed.ID,
			LastFetchedAt: sql.NullTime{Time: time.Now(), Valid: true}})
	if err != nil {
		return fmt.Errorf("failed to mark feed fetched: %w", err)
	}
	fmt.Printf("=== %v ===\n", feed.Channel.Title)
	for _, val := range feed.Channel.Item {
		fmt.Printf("- %v\n", val.Title)
		parsedTime, err := time.Parse(time.RFC1123Z, val.PubDate)
		if err != nil {
			return fmt.Errorf("failed to parse publication date: %w", err)
		}
		newPost := database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       val.Title,
			Url:         val.Link,
			Description: sql.NullString{String: val.Description, Valid: true},
			PublishedAt: parsedTime,
			FeedID:      nextFeed.ID,
		}
		_, err = s.db.CreatePost(context.Background(), newPost)
		if err != nil {
			errs := fmt.Sprintf("%v", err)
			if strings.Contains(strings.ToLower(errs), "url") {
				continue
			} else {
				return fmt.Errorf("error creating post: %w", err)
			}
		}

	}
	return nil
}
