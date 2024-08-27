package article

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	"time"
)

type ArticlePostgresStorage struct {
	db *sqlx.DB
}

func NewArticlePostgresStorage(db *sqlx.DB) *ArticlePostgresStorage {
	return &ArticlePostgresStorage{db: db}
}

type dbArticle struct {
	ID          int64          `db:"id"`
	SourceID    int64          `db:"id"`
	Title       string         `db:"title"`
	Link        string         `db:"link"`
	Summary     sql.NullString `db:"summary"`
	PublishedAt time.Time      `db:"published_at"`
	PostedAt    sql.NullTime   `db:"posted_at"`
	CreatedAt   time.Time      `db:"created_at"`
}
