package event

import (
	"context"

	"github.com/DeanXu2357/mychatbot/service/workout"
)

type event struct {
}

func (e *event) Events(ctx context.Context, uid string) ([]workout.Event, error) {
	//TODO implement me
	panic("implement me")
}

func (e *event) Create(ctx context.Context, uid string, name string, tags []workout.Tag) (workout.Event, error) {
	//TODO implement me
	panic("implement me")
}

func (e *event) Delete(ctx context.Context, eventID string) error {
	//TODO implement me
	panic("implement me")
}

func NewEditor() workout.EventEditor {
	return &event{}
}
