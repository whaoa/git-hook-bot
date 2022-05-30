package middlewares

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

func Logger(logger zerolog.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			err := next(c)
			request := c.Request()
			logger.Info().
				Str("method", request.Method).
				Str("url", request.URL.Path).
				Int("status", c.Response().Status).
				Int64("size", c.Response().Size).
				Str("usage", time.Now().Sub(start).String()).
				Msg("route logger")
			return err
		}
	}
}
