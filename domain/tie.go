package domain

import "time"

type Tie struct {
	Name      string `json:"name" db:"name"`
	TargetURL string `json:"target_url" db:"target_url"`

	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
