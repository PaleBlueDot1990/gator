package rssfeed

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/PaleBlueDot1990/gator/internal/config"
	"github.com/PaleBlueDot1990/gator/internal/database"
	"github.com/google/uuid"
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

	for _, item := range rssFeed.Channel.Item {
		publishedAt := sql.NullTime{}
		if t, err := time.Parse(time.RFC1123Z, item.PubDate); err == nil {
			publishedAt = sql.NullTime{
				Time:  t,
				Valid: true,
			}
		}

		dbQueryArgs := database.CreatePostParams {
			ID: uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Title: item.Title,
			Url: item.Link,
			Description:sql.NullString{String: item.Description, Valid: item.Description == ""},
			PublishedAt: publishedAt,
			FeedID: feed.ID,
		}

		_, err = state.DbQueries.CreatePost(ctx, dbQueryArgs)
		if err != nil {
			/*
			The error that a duplicate post is present in db will happen a lot
			because scraping job will repetitively fetch same rss feed items,
			and hence same rss feed links. And out posts table's URL column has 
			the UNIQUE constraint set on it.  
			*/
			if !strings.Contains(err.Error(), "violates unique constraint") {
				fmt.Printf("Error putting the rssfeed item in out database - %v", err)
			}
			continue 
		}
	}

	return nil
}