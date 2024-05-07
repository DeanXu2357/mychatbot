package logger

import (
	"context"

	"github.com/DeanXu2357/mychatbot/service/workout"
)

type workoutLogger struct{}

func (w *workoutLogger) Log(ctx context.Context, uid string, eventID string, rep int, weight int) error {
	//TODO implement me
	panic("implement me")
}

func (w *workoutLogger) Records(ctx context.Context, uid string, opts ...workout.QueryOption) ([]workout.Record, error) {
	//TODO implement me
	panic("implement me")
}

func (w *workoutLogger) Delete(ctx context.Context, uid string, recordID string) error {
	//TODO implement me
	panic("implement me")
}

func NewLogger() workout.Logger {
	return &workoutLogger{}
}
