package core

import (
	"net/http"

	"github.com/foxinuni/citas/core/controllers"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
)

type ListenAddress string

type SistemaCitas struct {
	address         ListenAddress
	citasController *controllers.CitaController
}

func NewSistemaCitas(address ListenAddress, citasController *controllers.CitaController) *SistemaCitas {
	return &SistemaCitas{
		address:         address,
		citasController: citasController,
	}
}

func (s *SistemaCitas) Listen() error {
	router := echo.New()
	router.HideBanner = true
	router.HidePort = true
	router.Validator = &CustomValidator{validator: validator.New()}

	router.GET("/", func(c echo.Context) error {
		return c.Redirect(http.StatusTemporaryRedirect, "/citas")
	})

	router.GET("/citas", s.citasController.GetAll)
	router.POST("/citas", s.citasController.Create)
	// router.GET("/citas/:id", s.citasController.GetById)
	// router.PUT("/citas/:id", s.citasController.Update)
	// router.DELETE("/citas/:id", s.citasController.Delete)

	router.GET("/static/*", echo.WrapHandler(http.StripPrefix("/static/", http.FileServer(http.Dir("static")))))

	return router.Start(string(s.address))
}

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}
