package config

import (
	"fmt"
	"reflect"

	"github.com/caarlos0/env"
)

// Config アプリケーション設定を表す構造体。
type Config struct {
	Port      int    `env:"PORT" envDefault:"50051"`
	Env       string `env:"GO_ENV" envDefault:"development"`
	LocalRoot string `env:"LOCAL_ROOT_PATH,required"`
	Debug     bool   `env:"DEBUG"`
}

var parseFuncMap = map[reflect.Type]env.ParserFunc{
	reflect.TypeOf(FakeTimeMode{}): parseFakeTimeMode,
}

// Get 環境変数からアプリケーション設定を取得する。
func Get() (*Config, error) {
	cfg := &Config{}
	if err := env.ParseWithFuncs(cfg, parseFuncMap); err != nil {
		return nil, fmt.Errorf("parse env: %w", err)
	}

	return cfg, nil
}
