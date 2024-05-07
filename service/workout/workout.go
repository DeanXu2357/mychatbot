package workout

import (
	"context"
	"time"
)

type Logger interface {
	Log(ctx context.Context, uid string, eventID string, rep int, weight int) error
	Records(ctx context.Context, uid string, opts ...QueryOption) ([]Record, error)
	Delete(ctx context.Context, uid string, recordID string) error
}

type EventEditor interface {
	Events(ctx context.Context, uid string) ([]Event, error)
	Create(ctx context.Context, uid string, name string, tags []Tag) (Event, error)
	Delete(ctx context.Context, eventID string) error
}

type QueryOption func(q *QueryFilter)

type QueryFilter struct {
	AfterID         string
	BeforeID        string
	BeforeCreatedAt time.Time
	AfterCreatedAt  time.Time
	Limit           int
}
