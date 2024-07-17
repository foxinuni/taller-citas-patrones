package controllers

import (
	"context"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

func RenderComponent(ctx echo.Context, status int, component templ.Component) error {
	buffer := templ.GetBuffer()
	defer templ.ReleaseBuffer(buffer)

	if err := component.Render(context.Background(), buffer); err != nil {
		return ctx.String(status, err.Error())
	}

	return ctx.HTML(status, buffer.String())
}
