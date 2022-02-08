package config

import "github.com/spf13/viper"

type EnvOptions struct {
	FileName  string
	FileType  string
	FilePath  string
	EnvPrefix string
}

type OptionFunc func(opts *EnvOptions)

func WithName(name string) OptionFunc {
	return func(opts *EnvOptions) {
		opts.FileName = name
	}
}

func WithType(t string) OptionFunc {
	return func(opts *EnvOptions) {
		opts.FileType = t
	}
}

func WithPath(path string) OptionFunc {
	return func(opts *EnvOptions) {
		opts.FilePath = path
	}
}

func WithPrefix(prefix string) OptionFunc {
	return func(opts *EnvOptions) {
		opts.EnvPrefix = prefix
	}
}

func NewEnvViper() (*viper.Viper, error) {
	d := EnvOptions{
		FileName:  ".env",
		FilePath:  ".",
		FileType:  "env",
		EnvPrefix: "zzn",
	}
	v := viper.New()
	// 2. 配置类型，支持 "json", "toml", "yaml", "yml", "properties",
	//             "props", "prop", "env", "dotenv"
	v.SetConfigType(d.FileType)
	// 3. 环境变量配置文件查找的路径，相对于 main.go
	v.AddConfigPath(d.FilePath)
	// 4. 设置环境变量前缀，用以区分 Go 的系统环境变量
	v.SetEnvPrefix(d.EnvPrefix)
	// 5. 读取环境变量（支持 flags）
	v.AutomaticEnv()
	// 6. 加载 env
	v.SetConfigName(d.FileName)
	// 7. 读取
	err := v.ReadInConfig()
	if err != nil {
		return nil, err
	}
	return v, nil
}
