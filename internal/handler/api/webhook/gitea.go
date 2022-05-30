package webhook

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"

	"github.com/whaoa/git-hook-bot/pkg/server"
	"github.com/whaoa/git-hook-bot/pkg/webhook/endpoint"
	"github.com/whaoa/git-hook-bot/pkg/webhook/git"
)

func Gitea(c echo.Context) error {
	endpointType := c.QueryParam("endpoint")
	secret := c.QueryParam("secret")

	if endpointType == "" || secret == "" {
		err := errors.New("miss required params: endpoint or secret.")
		return server.Resp(c, http.StatusBadRequest, "miss required params", err)
	}

	hook := git.ParseWebHook(git.Raw{
		Type: c.Request().Header.Get("X-GitHub-Event"),
		Body: c.Request().Body,
	})

	ep, err := endpoint.Create(endpointType, secret)
	if err != nil {
		return server.Resp(c, http.StatusBadRequest, "invalid endpoint", err)
	}

	if err = ep.SendMessage(hook); err != nil {
		return server.Resp(c, http.StatusInternalServerError, "send message to endpoint failed", err)
	}

	return server.Resp(c, http.StatusOK, "", nil)
}
