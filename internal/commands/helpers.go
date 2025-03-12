package commands

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/AhmettCelik/blog-aggregator/internal/database"
	"github.com/AhmettCelik/blog-aggregator/internal/rss"
	"github.com/AhmettCelik/blog-aggregator/internal/structure"
)

func scrapeFeeds(s *structure.State) error {
	feed, err := s.Database.GetNextFeedToFetch(context.Background())
	if err != nil {
		return fmt.Errorf("Error getting next feed to fetch: %v", err)
	}

	s.Database.MarkFeedFetched(context.Background(), feed.ID)

	fetchedFeedItems, err := rss.FetchFeed(context.Background(), feed.Url)
	if err != nil {
		return err
	}

	now := time.Now()

	for _, rssItem := range fetchedFeedItems.Channel.Item {
		parsedPubDate, err := time.Parse(time.RFC1123Z, rssItem.PubDate)
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("Error date could not be parsed: %v", err)
		}

		createPostParams := database.CreatePostParams{
			CreatedAt:   now,
			UpdatedAt:   now,
			Title:       rssItem.Title,
			Url:         rssItem.Link,
			Description: rssItem.Description,
			PublishedAt: parsedPubDate,
			FeedID:      feed.ID,
		}
		_, err = s.Database.CreatePost(context.Background(), createPostParams)
		if err != nil {
			if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "already exists") {
				fmt.Printf("Duplicated item detected: %s\n", rssItem.Title)
				continue
			}
			return fmt.Errorf("Error creating post: %v", err)
		}
	}
	return nil
}
