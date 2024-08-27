package storage

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/samber/lo"
	"news-parser-tg/internal/model"
	"time"
)

var (
	getAllSourcesQuery = `select * from sources`
	getSourceByIdQuery = `select * from sources 
                                  where source_id = $1`
	addSourceQuery = `insert into sources (source_name, feed_url, created_at) 
						   values ($1, $2, $3) 
					    returning source_id`
	deleteSourceByIdQuery = `delete from sources 
                                  where source_id = $1`
)

type SourcePostgresStorage struct {
	db *sqlx.DB
}

func (s *SourcePostgresStorage) Sources(ctx context.Context) ([]model.Source, error) {
	connx, err := s.db.Connx(ctx)
	if err != nil {
		return nil, err
	}
	defer connx.Close()

	var sources []dbSource
	if err := connx.SelectContext(ctx, &sources, getAllSourcesQuery); err != nil {
		return nil, err
	}

	return lo.Map(sources, func(source dbSource, _ int) model.Source {
		return model.Source(source)
	}), nil
}

func (s *SourcePostgresStorage) SourceByID(ctx context.Context, id int64) (*model.Source, error) {
	connx, err := s.db.Connx(ctx)
	if err != nil {
		return nil, err
	}
	defer connx.Close()

	var source dbSource
	if err := connx.GetContext(ctx, &source, getSourceByIdQuery, id); err != nil {
		return nil, err
	}

	return (*model.Source)(&source), nil
}

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

func (s *SourcePostgresStorage) Delete(ctx context.Context, id int64) error {
	connx, err := s.db.Connx(ctx)
	if err != nil {
		return err
	}
	defer connx.Close()

	if _, err := connx.ExecContext(ctx, deleteSourceByIdQuery, id); err != nil {
		return err
	}

	return nil
}

type dbSource struct {
	ID        int64     `db:"id"`
	Name      string    `db:"name"`
	FeedURL   string    `db:"feed_url"`
	CreatedAt time.Time `db:"created_at"`
}
