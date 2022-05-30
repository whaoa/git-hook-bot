package endpoint

import (
	"strings"

	"github.com/pkg/errors"

	"github.com/whaoa/git-hook-bot/pkg/webhook/git"
)

const (
	NameLarkSimpleBot = "lark-simple-bot"
)

var (
	ErrorEndPointNotSupported = errors.New("endpoint not supported")
)

type (
	client interface {
		name() string
		send(content string) error
		generateContent(hook *git.WebHookEvent) (string, error)
	}

	EndPoint struct {
		client client
	}
)

func (e EndPoint) SendMessage(hook *git.WebHookEvent) error {
	content, err := e.client.generateContent(hook)
	if err != nil {
		return err
	}
	return e.client.send(content)
}

func Create(endpoint, secret string) (*EndPoint, error) {
	switch strings.ToLower(endpoint) {
	case NameLarkSimpleBot:
		return &EndPoint{client: &LarkSimpleBot{secret: secret}}, nil
	}
	return nil, ErrorEndPointNotSupported
}
