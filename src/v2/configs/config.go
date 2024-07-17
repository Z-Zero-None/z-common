package configs

import (
	"fmt"
	"github.com/spf13/viper"
	"log/slog"
	"path/filepath"
)

func NewViperByFile(file string) (*viper.Viper, error) {
	cfg := viper.New()
	cfg.SetConfigFile(file)
	if err := cfg.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("apps.NewViper %s", err)
	}
	dir := filepath.Dir(file)
	for _, cfgFile := range cfg.GetStringSlice("includes") {
		v, err := NewViperByFile(filepath.Join(dir, cfgFile))
		if err != nil {
			slog.Error("apps.NewViper", "error", err)
			continue
		}

		cfg.MergeConfigMap(v.AllSettings())
	}
	return cfg, nil
}
