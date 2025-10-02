package rssfeed

import (
	"context"
	"encoding/xml"
	"html"
	"io"
	"net/http"
)

func FetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, feedURL, nil)
	if err != nil {
		return nil, err 
	}

	req.Header.Set("User-Agent", "gator")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err 
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err 
	}

	feed := &RSSFeed{}
	err = xml.Unmarshal(data, feed)
	if err != nil {
		return nil, err
	}

	decodeEscapedHTMLEntitities(feed)
	return feed, nil 
}

func decodeEscapedHTMLEntitities(feed *RSSFeed) {
	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)
	for idx := range feed.Channel.Item {
		feed.Channel.Item[idx].Title = 
			html.UnescapeString(feed.Channel.Item[idx].Title)
		feed.Channel.Item[idx].Description = 
			html.UnescapeString(feed.Channel.Item[idx].Description)
	}
}