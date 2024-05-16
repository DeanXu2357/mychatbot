package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"

	"github.com/DeanXu2357/mychatbot/service/workout"
)

type Event struct {
	ID        int64          `gorm:"autoIncrement;column:id"`
	Name      string         `gorm:"column:name"`
	UserID    string         `gorm:"column:user_id"`
	Tags      pq.StringArray `gorm:"column:tags;type:text[]"`
	CreatedAt time.Time      `gorm:"column:created_at"`
}

type event struct {
	db *gorm.DB
}

func (e *event) Events(ctx context.Context, uid string) ([]workout.Event, error) {
	var events []Event

	if err := e.db.WithContext(ctx).
		Table(TableEvent).
		Where("user_id", uid).
		Scan(&events).Error; err != nil {
		return nil, fmt.Errorf("failed to get events: %w", err)
	}

	output := make([]workout.Event, len(events))
	for _, v := range events {
		var tags []workout.Tag
		v.Tags.Scan(&tags)

		output = append(output, workout.Event{
			ID:        v.ID,
			Name:      v.Name,
			UserID:    v.UserID,
			Tags:      tags,
			CreatedAt: v.CreatedAt,
		})
	}

	return output, nil
}

func (e *event) Create(ctx context.Context, uid string, name string, tags []string) (workout.Event, error) {
	if err := e.db.WithContext(ctx).
		Table(TableEvent).
		Create(&Event{
			Name:      name,
			UserID:    uid,
			Tags:      tags,
			CreatedAt: time.Now(),
		}).Error; err != nil {
		return workout.Event{}, fmt.Errorf("failed to create event: %w", err)
	}

	return workout.Event{}, nil
}

func (e *event) Delete(ctx context.Context, eventID int64) error {
	if err := e.db.WithContext(ctx).Table(TableEvent).Where("id", eventID).Delete(&Event{}).Error; err != nil {
		return fmt.Errorf("failed to delete event: %w", err)
	}

	return nil
}

func NewEventEditor(db *gorm.DB) workout.EventEditor {
	return &event{db: db}
}
