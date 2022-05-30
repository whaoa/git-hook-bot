package global

import (
	"github.com/rs/zerolog"

	"github.com/whaoa/git-hook-bot/pkg/logger"
	"github.com/whaoa/git-hook-bot/pkg/settings"
)

var (
	AppSettings settings.App
	LogSettings settings.Log
	Settings    *settings.Settings
	Logger      zerolog.Logger
)

func Init() {
	Settings = settings.Setup()
	AppSettings = Settings.App
	LogSettings = Settings.Log

	Logger = logger.New(LogSettings.Level, LogSettings.Timestamp)
}
