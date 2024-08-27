package source

import (
	"context"
	"news-parser-tg/internal/model"
)

var (
	getSourceByIdQuery = `select * from sources 
                                  where source_id = $1`
)

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
