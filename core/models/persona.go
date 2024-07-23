package models

type Persona struct {
	Nombre   string `json:"nombre" form:"nombre" validate:"required"`
	Apellido string `json:"apellido" form:"apellido" validate:"required"`
	Cedula   string `json:"cedula" form:"cedula" validate:"required,numeric"`
	Edad     int    `json:"edad" form:"edad" validate:"required"`
}
