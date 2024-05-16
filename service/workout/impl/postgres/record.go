package postgres

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/DeanXu2357/mychatbot/service/workout"
)

const (
	PostgresImage    = ""
	PostgresImageTag = ""

	TableRecord = "records"
	TableEvent  = "events"
)

type Record struct {
	ID        int64     `gorm:"autoIncrement;column:id"`
	UserID    string    `gorm:"column:user_id"`
	Reps      int       `gorm:"column:reps"`
	Weight    int       `gorm:"column:weight"`
	CreatedAt time.Time `gorm:"column:created_at"`
	EventID   int64     `gorm:"column:event_id"`
}

type workoutRecord struct {
	db *gorm.DB
}

func (w *workoutRecord) Add(ctx context.Context, uid string, eventID int64, rep int, weight int) error {
	if err := w.db.WithContext(ctx).Table(TableRecord).Create(&Record{
		UserID:    uid,
		Reps:      rep,
		Weight:    weight,
		CreatedAt: time.Now(),
		EventID:   eventID,
	}).Error; err != nil {
		return fmt.Errorf("failed to create record: %w", err)
	}

	return nil
}

func (w *workoutRecord) Records(ctx context.Context, uid string, opts ...workout.QueryOption) ([]workout.Record, error) {
	queryfilter := workout.QueryFilter{
		Limit: 10,
	}

	for _, opt := range opts {
		opt(&queryfilter)
	}

	tx := w.db.WithContext(ctx).
		Table(TableRecord).
		Where("user_id", uid)

	if !queryfilter.AfterCreatedAt.IsZero() {
		tx = tx.Where("created_at > ?", queryfilter.AfterCreatedAt)
	}

	if !queryfilter.BeforeCreatedAt.IsZero() {
		tx = tx.Where("created_at < ?", queryfilter.BeforeCreatedAt)
	}

	var records []Record
	if err := tx.Limit(queryfilter.Limit).Find(&records).Error; err != nil {
		return nil, fmt.Errorf("failed to get records: %w", err)
	}

	var result []workout.Record
	for _, r := range records {
		result = append(result, workout.Record{
			ID:        r.ID,
			UserID:    r.UserID,
			Reps:      r.Reps,
			Weight:    r.Weight,
			CreatedAt: r.CreatedAt,
			Event: workout.Event{
				ID: r.EventID,
			},
		})
	}

	return result, nil
}

func (w *workoutRecord) Delete(ctx context.Context, uid string, recordID string) error {
	if err := w.db.WithContext(ctx).Table(TableRecord).
		Where("id", recordID).
		Where("user_id", uid).
		Delete(&Record{}).Error; err != nil {
		return fmt.Errorf("failed to delete record: %w", err)
	}

	return nil
}

func NewRecordEditor(db *gorm.DB) workout.RecordEditor {
	return &workoutRecord{db: db}
}

func WithRecordLimit(limit int) workout.QueryOption {
	return func(o *workout.QueryFilter) {
		o.Limit = limit
	}
}

func WithRecordBeforeCreatedAt(t time.Time) workout.QueryOption {
	return func(o *workout.QueryFilter) {
		o.BeforeCreatedAt = t
	}
}

func WithRecordAfterCreatedAt(t time.Time) workout.QueryOption {
	return func(o *workout.QueryFilter) {
		o.AfterCreatedAt = t
	}
}

func WithRecordAfterID(id string) workout.QueryOption {
	return func(o *workout.QueryFilter) {
		o.AfterID = id
	}
}

func WithRecordBeforeID(id string) workout.QueryOption {
	return func(o *workout.QueryFilter) {
		o.BeforeID = id
	}
}
