package article

import (
	"context"
	"github.com/samber/lo"
	"news-parser-tg/internal/model"
	"time"
)

var (
	getAllNonPostedArticlesQuery = `
  select * 
    from articles a
   where a.posted_at is null
     and a.published_at >= $1::timestamp
order by a.published_at desc
   limit $2
`
)

func (s *ArticlePostgresStorage) AllNonPosted(ctx context.Context, since time.Time, limit uint64) ([]model.Article, error) {
	connx, err := s.db.Connx(ctx)
	if err != nil {
		return nil, err
	}
	defer connx.Close()

	var articles []dbArticle
	if err := connx.SelectContext(
		ctx,
		&articles,
		getAllNonPostedArticlesQuery,
		since.UTC().Format(time.RFC3339),
		limit,
	); err != nil {
		return nil, err
	}

	return lo.Map(articles, func(article dbArticle, _ int) model.Article {
		return model.Article{
			ID:          article.ID,
			SourceID:    article.SourceID,
			Title:       article.Title,
			Link:        article.Link,
			Summary:     article.Summary.String,
			PublishedAt: article.PublishedAt,
			CreatedAt:   article.CreatedAt,
		}
	}), nil
}
