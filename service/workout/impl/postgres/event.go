package postgres

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/DeanXu2357/mychatbot/service/workout"
)

type Event struct {
	ID        int64         `gorm:"autoIncrement;column:id"`
	Name      string        `gorm:"column:name"`
	UserID    string        `gorm:"column:user_id"`
	Tags      []workout.Tag `gorm:"column:tags"`
	CreatedAt time.Time     `gorm:"column:created_at"`
}

type event struct {
	db *gorm.DB
}

func (e *event) Events(ctx context.Context, uid string) ([]workout.Event, error) {
	var events []workout.Event

	if err := e.db.WithContext(ctx).Table(TableEvent).Where("user_id", uid).Find(&events).Error; err != nil {
		return nil, fmt.Errorf("failed to get events: %w", err)
	}

	return events, nil
}

func (e *event) Create(ctx context.Context, uid string, name string, tags []workout.Tag) (workout.Event, error) {
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
