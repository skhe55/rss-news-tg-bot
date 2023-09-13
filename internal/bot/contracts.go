package bot

import (
	"context"
	"news-feed-bot/internal/model"
)

type SourceStorage interface {
	Sources(ctx context.Context) ([]model.Source, error)
	Add(ctx context.Context, source model.Source) (int64, error)
	Delete(ctx context.Context, id int64) error
	Update(ctx context.Context, source model.UpdateSource) error
}
