package main

import (
	"github.com/whaoa/git-hook-bot/internal/handler"
	"github.com/whaoa/git-hook-bot/pkg/global"
	"github.com/whaoa/git-hook-bot/pkg/server"
)

func main() {
	global.Init()

	s := server.Create(global.AppSettings.Mode, global.AppSettings.Address, global.Logger)
	handler.RegisterHandler(s)

	global.Logger.Info().
		Str("mode", global.AppSettings.Mode).
		Str("address", global.AppSettings.Address).
		Msg("app start")

	if err := s.Start(); err != nil {
		global.Logger.Error().Err(err).Msg("error in app start")
		panic(err)
	}
}
