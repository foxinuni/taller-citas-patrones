package stores

import (
	"context"
	"errors"

	"github.com/foxinuni/citas/core/models"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

var _ CitaStore = &PostgresCitaStore{}

type PostgresCitaStore struct {
	pool *pgxpool.Pool
}

func NewPostgresCitaStore(pool *pgxpool.Pool) *PostgresCitaStore {
	return &PostgresCitaStore{pool: pool}
}

func (s *PostgresCitaStore) GetAll(filter CitaStoreFilter) ([]models.Cita, error) {
	// Check if limit is 0
	if filter.Limit == 0 {
		filter.Limit = 10
	}

	var citas []models.Cita
	rows, err := s.pool.Query(context.Background(), `
		SELECT id, nombre, apellido, cedula, edad, fecha
		FROM citas
		WHERE fecha::date = $1
		ORDER BY fecha DESC
		LIMIT $2 OFFSET $3
	`, filter.Date, filter.Limit, filter.Limit*(filter.Page-1))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Iterate over the rows
	for rows.Next() {
		var cita models.Cita
		err := rows.Scan(&cita.ID, &cita.Persona.Nombre, &cita.Persona.Apellido, &cita.Persona.Cedula, &cita.Persona.Edad, &cita.Fecha)
		if err != nil {
			return nil, err
		}
		citas = append(citas, cita)
	}

	return citas, nil
}

func (s *PostgresCitaStore) GetById(id string) (models.Cita, error) {
	var cita models.Cita
	err := s.pool.QueryRow(context.Background(), `
		SELECT id, nombre, apellido, cedula, edad, fecha
		FROM citas
		WHERE id = $1
	`, id).Scan(&cita.ID, &cita.Persona.Nombre, &cita.Persona.Apellido, &cita.Persona.Cedula, &cita.Persona.Edad, &cita.Fecha)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Cita{}, ErrCitaNotFound
		}

		return models.Cita{}, err
	}

	return cita, nil
}

func (s *PostgresCitaStore) Create(cita *models.Cita) error {
	obj, err := s.pool.Exec(context.Background(), `
		INSERT INTO citas (id, nombre, apellido, cedula, edad, fecha)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, cita.ID, cita.Persona.Nombre, cita.Persona.Apellido, cita.Persona.Cedula, cita.Persona.Edad, cita.Fecha)
	if err != nil {
		return err
	}

	if obj.RowsAffected() == 0 {
		return ErrCitaExists
	}

	return nil
}

func (s *PostgresCitaStore) Update(cita *models.Cita) error {
	obj, err := s.pool.Exec(context.Background(), `
		UPDATE citas
		SET 
			nombre = CASE WHEN $1 != '' THEN $1 ELSE nombre END,
			apellido = CASE WHEN $2 != '' THEN $2 ELSE apellido END,
			cedula = CASE WHEN $3 != '' THEN $3 ELSE cedula END,
			edad = CASE WHEN $4 != 0 THEN $4 ELSE edad END,
			fecha = CASE WHEN $5 != '' THEN $5 ELSE fecha END
		WHERE id = $6
	`, cita.Persona.Nombre, cita.Persona.Apellido, cita.Persona.Cedula, cita.Persona.Edad, cita.Fecha, cita.ID)
	if err != nil {
		return err
	}

	if obj.RowsAffected() == 0 {
		return ErrCitaNotFound
	}

	return nil
}

func (s *PostgresCitaStore) Delete(cita *models.Cita) error {
	obj, err := s.pool.Exec(context.Background(), `
		DELETE FROM citas
		WHERE id = $1
	`, cita.ID)
	if err != nil {
		return err
	}

	if obj.RowsAffected() == 0 {
		return ErrCitaNotFound
	}

	return nil
}
