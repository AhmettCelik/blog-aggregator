package commands

import (
	"context"
	"fmt"

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

	for _, rssItem := range fetchedFeedItems.Channel.Item {
		fmt.Printf("%s\n", rssItem.Title)
	}
	return nil
}
