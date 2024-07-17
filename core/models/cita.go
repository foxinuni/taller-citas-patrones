package models

import (
	"time"
)

type Cita struct {
	ID      string    `json:"id"`
	Persona Persona   `json:"persona" validate:"required"`
	Fecha   time.Time `json:"fecha" validate:"required"`
}
