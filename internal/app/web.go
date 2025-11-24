package app

import (
	"context"
	"errors"
	"time"

	"github.com/labstack/echo/v4"
)

type WebApp struct {
	echo        *echo.Echo
	restAddress string
}

func NewWebApp(restAddress string) *WebApp {
	return &WebApp{
		restAddress: restAddress,
	}
}

func (app *WebApp) Start() error {
	app.echo = echo.New()

	app.echo.Static("/", "frontend")

	if err := app.echo.Start(app.restAddress); err != nil {
		return errors.New("error starting web server: " + err.Error())
	}

	return nil
}

func (app *WebApp) GracefulShutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	if err := app.echo.Shutdown(ctx); err != nil {
		app.echo.Logger.Fatal(err)
	}
}
