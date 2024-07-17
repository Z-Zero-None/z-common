package apps

import (
	"github.com/spf13/viper"
	"log/slog"
	"os"
	"z-common/src/v2/configs"
)

func NewAppConfig(file string) (*viper.Viper, error) {
	cfg, err := configs.NewViperByFile(file)
	if err != nil {
		return nil, err
	}
	// 设置环境
	SetMode(cfg.GetString("env"))

	cfg.WatchConfig()

	namespace = cfg.GetString("service.name") + "."

	lvl := new(slog.LevelVar)
	lvl.UnmarshalText([]byte(cfg.GetString("log-level")))
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
		AddSource:   false,
		Level:       lvl,
		ReplaceAttr: nil,
	})))
	return cfg, nil
}
