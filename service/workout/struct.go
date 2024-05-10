package workout

import "time"

type Tag string

type Event struct {
	ID        int64
	Name      string
	UserID    string
	Tags      []Tag
	CreatedAt time.Time
}

type Record struct {
	ID        int64
	UserID    string
	Reps      int
	Weight    int
	CreatedAt time.Time
	Event
}
