package source

import "context"

var (
	deleteSourceByIdQuery = `
delete from sources 
      where source_id = $1
`
)

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
