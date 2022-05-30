package handler

import (
    "github.com/whaoa/git-hook-bot/internal/handler/api/webhook"
    "github.com/whaoa/git-hook-bot/pkg/server"
)

func RegisterHandler(server *server.Server) {
    wh := server.Group("/webhook")
    {
        wh.Any("/gitea", webhook.Gitea)
    }
}
