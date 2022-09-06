package domain

import "time"

type Status struct {
	Alive bool      `json:"alive"`
	Time  time.Time `json:"time"`
}

func CurrentStatus() *Status {
	return &Status{true, time.Now()}
}
