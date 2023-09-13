package storage

import (
	"context"
	"fmt"
	"news-feed-bot/internal/model"
	"news-feed-bot/internal/utils"
	"reflect"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

type SourcePostgresStorage struct {
	db *sqlx.DB
}

type dbSource struct {
	Id        int64     `db:"id"`
	Name      string    `db:"name"`
	FeedUrl   string    `db:"feed_url"`
	CreatedAt time.Time `db:"created_at"`
}

func NewArticleStorage(db *sqlx.DB) *ArticlePostgresStorage {
	return &ArticlePostgresStorage{db: db}
}

func (s *SourcePostgresStorage) Sources(ctx context.Context) ([]model.Source, error) {
	conn, err := s.db.Connx(ctx)

	if err != nil {
		return nil, err
	}
	defer conn.Close()

	var sources []dbSource
	if err := conn.SelectContext(ctx, &sources, `SELECT * FROM sources`); err != nil {
		return nil, err
	}

	return utils.Map(sources, func(source dbSource, _ int) model.Source {
		return model.Source{
			Id:        source.Id,
			Name:      source.Name,
			FeedUrl:   source.FeedUrl,
			CreatedAt: source.CreatedAt,
		}
	}), nil
}

func (s *SourcePostgresStorage) SourceById(ctx context.Context, id int64) (*model.Source, error) {
	conn, err := s.db.Connx(ctx)

	if err != nil {
		return nil, err
	}
	defer conn.Close()

	var source dbSource
	if err := conn.GetContext(ctx, &source, `SELECT * FROM sources where id = $1`, id); err != nil {
		return nil, err
	}

	return (*model.Source)(&source), nil
}

func (s *SourcePostgresStorage) Add(ctx context.Context, source model.Source) (int64, error) {
	conn, err := s.db.Connx(ctx)

	if err != nil {
		return 0, err
	}
	defer conn.Close()

	var id int64

	row := conn.QueryRowxContext(
		ctx,
		`INSERT INTO sources (name, feed_url, created_at) VALUES ($1, $2, $3) RETURNING id`,
		source.Name,
		source.FeedUrl,
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

func (s *SourcePostgresStorage) Update(ctx context.Context, source model.UpdateSource) error {
	conn, err := s.db.Connx(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	valueOfSource := reflect.ValueOf(source)
	typeOfSource := valueOfSource.Type()

	var fields []string

	for i := 1; i < valueOfSource.NumField(); i++ {
		if valueOfSource.Field(i).Interface() != "" && typeOfSource.Field(i).Name != "URL" {
			fields = append(fields, fmt.Sprintf("%v=:%v", strings.ToLower(typeOfSource.Field(i).Name), strings.ToLower(typeOfSource.Field(i).Name)))
		}
		if valueOfSource.Field(i).Interface() != "" && typeOfSource.Field(i).Name == "URL" {
			fields = append(fields, fmt.Sprintf("feed_url=:%v", strings.ToLower(typeOfSource.Field(i).Name)))
		}
	}

	namedQuery := fmt.Sprintf("UPDATE sources SET %v WHERE id=:id", strings.Join(fields, ", "))

	query, args, err := sqlx.Named(namedQuery, source)
	if err != nil {
		return err
	}

	if len(args) == 1 {
		return nil
	}

	if _, err := conn.ExecContext(ctx, conn.Rebind(query), args...); err != nil {
		return err
	}

	return nil
}

func (s *SourcePostgresStorage) Delete(ctx context.Context, id int64) error {
	conn, err := s.db.Connx(ctx)

	if err != nil {
		return err
	}
	defer conn.Close()

	if _, err := conn.ExecContext(ctx, `DELETE FROM sources WHERE id = $1`, id); err != nil {
		return err
	}

	return nil
}
