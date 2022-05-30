package server

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/whaoa/git-hook-bot/pkg/server/middlewares"
)

type Server struct {
	*echo.Echo
	mode    string
	address string
	logger  zerolog.Logger
}

func (s Server) Start() error {
	return errors.Wrap(s.Echo.Start(s.address), "echo server start error")
}

func Create(mode, address string, logger zerolog.Logger) *Server {
	engine := echo.New()

	engine.HidePort = true
	engine.HideBanner = true
	engine.Logger.SetLevel(log.OFF)

	engine.HTTPErrorHandler = func(err error, ctx echo.Context) {
		code := http.StatusInternalServerError
		message := http.StatusText(http.StatusInternalServerError)

		httpError, ok := err.(*echo.HTTPError)
		if ok {
			code = httpError.Code
			message = fmt.Sprintf("http error: %s", httpError.Message)
		}

		e := Resp(ctx, code, message, err)
		if e != nil {
			logger.Error().Err(e).AnErr("raw-errors", err).Msg("error in response method")
		}
	}

	engine.Use(
		middleware.RecoverWithConfig(middleware.RecoverConfig{
			LogErrorFunc: func(c echo.Context, err error, stack []byte) error {
				logger.Error().Err(err).Bytes("stack", stack).Msg("panic recovered")
				return nil
			},
		}),
		middlewares.Logger(logger),
	)

	return &Server{Echo: engine, mode: mode, address: address, logger: logger}
}
