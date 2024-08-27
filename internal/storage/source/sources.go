package source

import (
	"context"
	"github.com/samber/lo"
	"news-parser-tg/internal/model"
)

var (
	getAllSourcesQuery = `
select * 
  from sources
`
)

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
