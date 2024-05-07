package workout

import "time"

type Tag string

type Event struct {
	Name      string
	UserID    string
	Tags      []Tag
	CreatedAt time.Time
}

type Record struct {
	UserID    string
	Reps      int
	Weight    int
	CreatedAt time.Time
	Event
}
