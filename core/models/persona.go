package models

type Persona struct {
	Nombre   string `json:"nombre" param:"nombre" validate:"required"`
	Apellido string `json:"apellido" param:"apellido" validate:"required"`
	Cedula   string `json:"cedula" param:"cedula" validate:"required,numeric"`
	Edad     int    `json:"edad" param:"edad" validate:"required"`
}
