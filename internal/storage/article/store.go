package article

import (
	"context"
	"news-parser-tg/internal/model"
)

var (
	addArticleQuery = `
insert into articles(source_id, title, link, summary, published_at)
     values ($1, $2, $3, $4, $5)
on conflict do nothing
`
)

func (s *ArticlePostgresStorage) Store(ctx context.Context, article model.Article) error {
	connx, err := s.db.Connx(ctx)
	if err != nil {
		return err
	}
	defer connx.Close()

	if _, err := connx.ExecContext(
		ctx,
		addArticleQuery,
		article.SourceID,
		article.Title,
		article.Link,
		article.Summary,
		article.PublishedAt,
	); err != nil {
		return err
	}
}
