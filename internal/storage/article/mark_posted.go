package article

import (
	"context"
	"time"
)

var (
	markPostedArticleQuery = `
update articles
   set posted_at = $1::timestamp
 where id = $2
`
)

func (s *ArticlePostgresStorage) MarkPosted(ctx context.Context, id int64) error {
	connx, err := s.db.Connx(ctx)
	if err != nil {
		return err
	}
	defer connx.Close()

	if _, err := connx.ExecContext(
		ctx,
		markPostedArticleQuery,
		time.Now().UTC().Format(time.RFC3339),
		id); err != nil {
		return err
	}

	return nil
}
