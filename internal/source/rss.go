package source

import (
	"context"
	"news-feed-bot/internal/model"

	"github.com/SlyMarbo/rss"
)

type RSSSource struct {
	URL        string
	SourceId   int64
	SourceName string
}

func NewRSSSourceFromModel(m model.Source) RSSSource {
	return RSSSource{
		URL:        m.FeedUrl,
		SourceId:   m.Id,
		SourceName: m.Name,
	}
}

func (s RSSSource) Fetch(ctx context.Context) ([]model.Item, error) {
	feed, err := s.loadFeed(ctx, s.URL)
	mappedFeed := make([]model.Item, len(feed.Items))
	/* Вынести в отдельный метод Map + Utests */
	for _, v := range feed.Items {
		elem := model.Item{
			Title:      v.Title,
			Categories: v.Categories,
			Link:       v.Link,
			Date:       v.Date,
			Summary:    v.Summary,
			SourceName: s.SourceName,
		}
		mappedFeed = append(mappedFeed, elem)
	}

	if err != nil {
		return nil, err
	}

	return mappedFeed, nil
}

func (s RSSSource) loadFeed(ctx context.Context, url string) (*rss.Feed, error) {
	var (
		feedCh = make(chan *rss.Feed)
		errCh  = make(chan error)
	)

	go func() {
		feed, err := rss.Fetch(url)
		if err != nil {
			errCh <- err
			return
		}

		feedCh <- feed
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case err := <-errCh:
		return nil, err
	case feed := <-feedCh:
		return feed, nil
	}
}
