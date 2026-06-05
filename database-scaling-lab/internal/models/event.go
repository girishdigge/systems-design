package models

import (
	"time"

	"github.com/google/uuid"
)

type Event struct {
	UserID    uuid.UUID `json:"user_id"`
	EventType string    `json:"event_type"`
	CreatedAt time.Time `json:"created_at"`
}
