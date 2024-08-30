package stores

import (
	"errors"
	"time"

	"github.com/foxinuni/citas/core/models"
)

type CitaStoreFilter struct {
	Date  time.Time
	Limit int
	Page  int
}

var (
	ErrCitaNotFound = errors.New("failed to find cita with the specified ID")
	ErrInvalidId    = errors.New("invalid ID format")
	ErrCitaExists   = errors.New("cita already exists on this date")
)

type CitaStore interface {
	GetAll(filter CitaStoreFilter) ([]models.Cita, error)
	GetById(id string) (models.Cita, error)
	Create(cita *models.Cita) error
	Update(cita *models.Cita) error
	Delete(cita *models.Cita) error
}
