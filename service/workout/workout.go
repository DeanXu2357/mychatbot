package workout

import (
	"context"
	"time"
)

type RecordEditor interface {
	Add(ctx context.Context, uid string, eventID int64, rep int, weight int) error
	Records(ctx context.Context, uid string, opts ...QueryOption) ([]Record, error)
	Delete(ctx context.Context, uid string, recordID string) error
}

type EventEditor interface {
	Events(ctx context.Context, uid string) ([]Event, error)
	Create(ctx context.Context, uid string, name string, tags []string) (Event, error)
	Delete(ctx context.Context, eventID int64) error
}

type QueryOption func(q *QueryFilter)

type QueryFilter struct {
	AfterID         string
	BeforeID        string
	BeforeCreatedAt time.Time
	AfterCreatedAt  time.Time
	Limit           int
}
