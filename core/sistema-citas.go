package core

import (
	"net/http"

	"github.com/foxinuni/citas/core/controllers"
	"github.com/foxinuni/citas/core/views"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type SistemaCitas struct {
	address         string
	citasController *controllers.CitaController
}

func NewSistemaCitas(address string, citasController *controllers.CitaController) *SistemaCitas {
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

	router.Use(s.ErrorHandlerMiddleware)
	router.GET("/", func(c echo.Context) error {
		return c.Redirect(http.StatusTemporaryRedirect, "/citas")
	})

	router.GET("/citas", s.citasController.GetAll)
	router.POST("/citas", s.citasController.Create)
	router.GET("/citas/new", s.citasController.New)
	// router.GET("/citas/:id", s.citasController.GetById)
	// router.PUT("/citas/:id", s.citasController.Update)
	// router.DELETE("/citas/:id", s.citasController.Delete)

	router.GET("/static/*", echo.WrapHandler(http.StripPrefix("/static/", http.FileServer(http.Dir("static")))))

	return router.Start(string(s.address))
}

func (s *SistemaCitas) ErrorHandlerMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		log.Info().Msgf("%s %s from [%s]", c.Request().Method, c.Request().URL.Path, c.Request().RemoteAddr)

		if err := next(c); err != nil {
			log.Error().Err(err).Msg("Error processing request")

			var code int
			if he, ok := err.(*echo.HTTPError); ok {
				code = he.Code
			} else {
				code = http.StatusInternalServerError
			}

			view := views.ViewError(err)
			controllers.RenderComponent(c, code, view)
		}

		return nil
	}
}

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}
