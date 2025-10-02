package rssfeed

import (
	"context"
	"fmt"

	"github.com/PaleBlueDot1990/gator/internal/config"
)

func ScrapeFeeds(ctx context.Context, state *config.State) error {
	feed, err := state.DbQueries.GetNextFeedToFetch(ctx)
	if err != nil {
		return err
	}

	err = state.DbQueries.MarkFeedFetched(ctx, feed.ID)
	if err != nil {
		return err
	}

	rssFeed, err := FetchFeed(ctx, feed.Url)
	if err != nil {
		return err
	}

	fmt.Printf("------------------------------------------------------\n")
	fmt.Printf("Feeds scraped from - %s\n", feed.Url)
	for idx, item := range rssFeed.Channel.Item {
		fmt.Printf("%d. %s\n", idx, item.Title)
	}
	fmt.Printf("------------------------------------------------------\n")
	return nil
}