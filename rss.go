package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"time"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, feedURL, nil)
	if err != nil {
		return &RSSFeed{}, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("User-Agent", "gator")
	c := http.Client{
		Timeout: 10 * time.Second,
	}

	res, err := c.Do(req)
	if err != nil {
		return &RSSFeed{}, fmt.Errorf("request to server failed: %w", err)
	}
	defer res.Body.Close()
	rssfeed := &RSSFeed{}
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return &RSSFeed{}, fmt.Errorf("failed to read response: %w", err)
	}
	err = xml.Unmarshal(data, rssfeed)
	if err != nil {
		return &RSSFeed{}, fmt.Errorf("failed to unmarshal response: %w", err)
	}
	rssfeed.Channel.Title = html.UnescapeString(rssfeed.Channel.Title)
	rssfeed.Channel.Description = html.UnescapeString(rssfeed.Channel.Description)
	return rssfeed, nil
}
