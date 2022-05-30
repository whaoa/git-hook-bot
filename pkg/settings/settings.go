package settings

import (
	"strings"

	"github.com/spf13/viper"
)

const (
	ModeRelease     = "release"
	ModeDevelopment = "development"

	DefaultAppMode            = ModeDevelopment
	DefaultAppAddress         = ":2333"
	DefaultLogLevel           = "Debug"
	DefaultLogTimeStampFormat = "2006-01-02 15:04:05"
)

type (
	App struct {
		Mode    string `toml:"mode" json:"mode"`
		Address string `toml:"address" json:"address"`
	}
	Log struct {
		Level     string `toml:"level" json:"level"`
		Timestamp string `toml:"timestamp" json:"timestamp"`
	}
	Settings struct {
		App App `toml:"app" json:"app"`
		Log Log `toml:"log" json:"log"`
	}
)

var v = viper.New()

func init() {
	v.SetConfigName("config")
	v.SetConfigType("toml")

	v.AddConfigPath("./")
	v.AddConfigPath("./config")
	v.AddConfigPath("./temp")

	v.SetDefault("app.mode", DefaultAppMode)
	v.SetDefault("app.address", DefaultAppAddress)

	v.SetDefault("log.level", DefaultLogLevel)
	v.SetDefault("log.timestamp", DefaultLogTimeStampFormat)

	v.SetEnvKeyReplacer(strings.NewReplacer(".", "__"))
	v.AutomaticEnv()
}

func Setup() *Settings {
	settings := &Settings{}

	var err error

	if err = v.ReadInConfig(); err != nil {
		panic(err)
	}
	if err = v.Unmarshal(settings); err != nil {
		panic(err)
	}

	if settings.App.Mode != ModeRelease {
		settings.App.Mode = ModeDevelopment
	}

	return settings
}
