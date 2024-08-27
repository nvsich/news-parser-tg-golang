package source

import (
	"context"
	"news-parser-tg/internal/model"
)

var (
	addSourceQuery = `
insert into sources (source_name, feed_url, created_at) 
	 values ($1, $2, $3) 
  returning id
`
)

func (s *SourcePostgresStorage) Add(ctx context.Context, source *model.Source) (int64, error) {
	connx, err := s.db.Connx(ctx)
	if err != nil {
		return 0, err
	}
	defer connx.Close()

	var id int64
	row := connx.QueryRowxContext(
		ctx,
		addSourceQuery,
		source.Name,
		source.FeedURL,
		source.CreatedAt,
	)

	if err := row.Err(); err != nil {
		return 0, err
	}

	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}
