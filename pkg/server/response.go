package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Response struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

func Resp(ctx echo.Context, code int, message string, err error) error {
	if code == 0 {
		if err != nil {
			code = http.StatusInternalServerError
		} else {
			code = http.StatusOK
		}
	}
	if message == "" {
		message = "ok"
	}

	errorMessage := ""
	if err != nil {
		errorMessage = err.Error()
	}

	return ctx.JSON(code, Response{Message: message, Error: errorMessage})
}
