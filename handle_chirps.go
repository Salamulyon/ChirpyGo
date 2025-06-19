package main

import (
	"time"

	"github.com/google/uuid"
)

type Chirp struct {
	body       string    `json:"body"`
	user_id    uuid.UUID `json:"user_id"`
	created_at time.Time `json:"created_at"`
	updated_at time.Time `json:"updated_at"`
}
