package controllers

import (
	"errors"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/foxinuni/citas/core/models"
	"github.com/foxinuni/citas/core/stores"
	"github.com/foxinuni/citas/core/views"
	"github.com/labstack/echo/v4"
)

type CitaController struct {
	store stores.CitaStore
}

func NewCitaController(store stores.CitaStore) *CitaController {
	return &CitaController{store: store}
}

func (c *CitaController) GetAll(ctx echo.Context) error {
	filters := stores.CitaStoreFilter{}

	// Se obtiene el parametro date de la URL
	date := ctx.QueryParam("date")
	if date == "" {
		today := time.Now()
		filters.Date = time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, time.Local)
	} else {
		filters.Date, _ = time.Parse("2006-01-02", date)
	}

	// Se obtiene el parametro page de la URL
	limit := ctx.QueryParam("limit")
	if limit == "" {
		filters.Limit = 10
	} else {
		if v, err := strconv.Atoi(limit); err == nil {
			filters.Limit = v
		} else {
			return ctx.JSON(http.StatusBadRequest, "limit parameter must be a number")
		}
	}

	// Se obtiene el parametro page de la URL
	page := ctx.QueryParam("page")
	if page == "" {
		filters.Page = 1
	} else {
		if v, err := strconv.Atoi(page); err == nil {
			if v < 1 {
				return ctx.JSON(http.StatusBadRequest, "page parameter must be greater than 0")
			}

			filters.Page = v
		} else {
			return ctx.JSON(http.StatusBadRequest, "page parameter must be a number")
		}
	}

	// Se obtienen las citas
	citas, err := c.store.GetAll(filters)
	if err != nil {
		if errors.Is(err, stores.ErrCitaNotFound) {
			return ctx.JSON(http.StatusNotFound, err.Error())
		} else if errors.Is(err, stores.ErrInvalidId) {
			return ctx.JSON(http.StatusBadRequest, err.Error())
		} else if errors.Is(err, stores.ErrCitaExists) {
			return ctx.JSON(http.StatusConflict, err.Error())
		} else {
			return ctx.JSON(http.StatusInternalServerError, err.Error())
		}
	}

	// Se ordenan las citas por hora
	sort.Slice(citas, func(i, j int) bool {
		return citas[i].Fecha.Before(citas[j].Fecha)
	})

	// Se renderiza la vista
	view := views.ViewCitaList(citas)
	return RenderComponent(ctx, http.StatusOK, view)
}

func (c *CitaController) Create(ctx echo.Context) error {
	cita := new(models.Cita)
	if err := ctx.Bind(cita); err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	if err := ctx.Validate(cita); err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	if err := c.store.Create(cita); err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusCreated, cita)
}
