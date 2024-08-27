package fetcher

import (
	"context"
	"news-parser-tg/internal/model"
	"news-parser-tg/internal/source"
	"strings"
	"sync"
	"time"
)

type ArticleSaver interface {
	Store(ctx context.Context, article model.Article) error
}

type SourceProvider interface {
	Sources(ctx context.Context) ([]model.Source, error)
}

type Source interface {
	ID() int64
	Name() string
	Fetch(ctx context.Context) ([]model.Item, error)
}

type Fetcher struct {
	articleSaver   ArticleSaver
	sourceProvider SourceProvider

	fetchInterval  time.Duration
	filterKeywords []string
}

func New(
	articleSaver ArticleSaver,
	sourceProvider SourceProvider,
	fetchInterval time.Duration,
	filterKeywords []string,
) *Fetcher {
	return &Fetcher{
		articleSaver:   articleSaver,
		sourceProvider: sourceProvider,
		fetchInterval:  fetchInterval,
		filterKeywords: filterKeywords,
	}
}

func (f *Fetcher) Start(ctx context.Context) error {
	ticker := time.NewTicker(f.fetchInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if err := f.Fetch(ctx); err != nil {
				return err
			}
		}
	}
}

func (f *Fetcher) Fetch(ctx context.Context) error {
	sources, err := f.sourceProvider.Sources(ctx)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup

	for _, src := range sources {
		wg.Add(1)

		rssSource := source.NewRSSSourceFromModel(&src)

		go func(source Source) {
			defer wg.Done()

			items, err := source.Fetch(ctx)
			if err != nil {
				// TODO: log
				return
			}

			if err := f.processItems(ctx, source, items); err != nil {
				// TODO: log
				return
			}
		}(rssSource)
	}

	wg.Wait()

	return nil
}

func (f *Fetcher) processItems(ctx context.Context, source Source, items []model.Item) error {
	for _, item := range items {
		item.Date = time.Now().UTC()

		if f.itemShouldBeSkipped(item) {
			continue
		}

		if err := f.articleSaver.Store(ctx, model.Article{
			SourceID:    source.ID(),
			Title:       item.Title,
			Link:        item.Link,
			Summary:     item.Summary,
			PublishedAt: item.Date,
		}); err != nil {
			return err
		}
	}

	return nil
}

func (f *Fetcher) itemShouldBeSkipped(item model.Item) bool {
	categories := make([]string, len(item.Categories))
	copy(categories, item.Categories)
	categoriesSet := map[string]bool{}

	for _, category := range categories {
		categoriesSet[category] = true
	}

	for _, keyword := range f.filterKeywords {
		if categoriesSet[keyword] || strings.Contains(item.Title, keyword) {
			return true
		}
	}

	return false
}
