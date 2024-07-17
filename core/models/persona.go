package models

type Persona struct {
	Nombre   string `json:"nombre" validate:"required"`
	Apellido string `json:"apellido" validate:"required"`
	Cedula   string `json:"cedula" validate:"required,numeric"`
	Edad     int    `json:"edad" validate:"required"`
}
