package storage

import (
	"context"
	"database/sql"
	"news-feed-bot/internal/model"
	"news-feed-bot/internal/utils"
	"time"

	"github.com/jmoiron/sqlx"
)

type ArticlePostgresStorage struct {
	db *sqlx.DB
}

type dbArticle struct {
	Id          int64        `db:"id"`
	SourceId    int64        `db:"source_id"`
	Title       string       `db:"title"`
	Link        string       `db:"link"`
	Summary     string       `db:"summary"`
	PublishedAt time.Time    `db:"published_at"`
	PostedAt    sql.NullTime `db:"posted_at"`
	CreatedAt   time.Time    `db:"created_at"`
}

func NewSourceStorage(db *sqlx.DB) *SourcePostgresStorage {
	return &SourcePostgresStorage{db: db}
}

func (a *ArticlePostgresStorage) Articles(ctx context.Context) ([]model.Article, error) {
	conn, err := a.db.Connx(ctx)

	if err != nil {
		return nil, err
	}

	defer conn.Close()

	var articles []dbArticle
	if err := conn.SelectContext(ctx, &articles, `SELECT * FROM articles`); err != nil {
		return nil, err
	}

	return utils.Map(articles, func(article dbArticle, _ int) model.Article {
		return model.Article{
			Id:          article.Id,
			SourceId:    int64(article.SourceId),
			Title:       article.Title,
			Link:        article.Link,
			Summary:     article.Summary,
			PublishedAt: article.PublishedAt,
			PostedAt:    article.PostedAt.Time,
			CreatedAt:   article.CreatedAt,
		}
	}), nil
}

func (a *ArticlePostgresStorage) Store(ctx context.Context, article model.Article) error {
	conn, err := a.db.Connx(ctx)
	if err != nil {
		return err
	}

	defer conn.Close()

	if _, err := conn.ExecContext(
		ctx,
		`INSERT INTO articles (source_id, title, link, summary, published_at) VALUES ($1, $2, $3, $4, $5) ON CONFLICT DO NOTHING;`,
		article.SourceId,
		article.Title,
		article.Link,
		article.Summary,
		article.PublishedAt,
	); err != nil {
		return err
	}

	return nil
}

func (a *ArticlePostgresStorage) AllNotPosted(ctx context.Context, since time.Time, limit uint64) ([]model.Article, error) {
	conn, err := a.db.Connx(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	var articles []dbArticle
	if err := conn.SelectContext(
		ctx,
		&articles,
		`SELECT * FROM articles WHERE posted_at IS NULL AND published_at <= $1::timestamp ORDER BY published_at DESC LIMIT $2`,
		since.UTC().Format(time.RFC3339), limit,
	); err != nil {
		return nil, err
	}

	return utils.Map(articles, func(article dbArticle, _ int) model.Article {
		return model.Article{
			Id:          article.Id,
			SourceId:    article.SourceId,
			Title:       article.Title,
			Link:        article.Link,
			Summary:     article.Summary,
			PostedAt:    article.PostedAt.Time,
			PublishedAt: article.PublishedAt,
			CreatedAt:   article.CreatedAt,
		}
	}), nil
}

func (a *ArticlePostgresStorage) MarkPosted(ctx context.Context, article model.Article) error {
	conn, err := a.db.Connx(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	if _, err := conn.ExecContext(
		ctx,
		`UPDATE articles SET posted_at = $1::timestamp WHERE id = $2`,
		time.Now().UTC().Format(time.RFC3339),
		article.Id,
	); err != nil {
		return err
	}

	return nil
}
