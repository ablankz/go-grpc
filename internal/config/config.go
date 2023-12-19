package config

import (
	"fmt"
	"reflect"

	"github.com/caarlos0/env"
)

// Config アプリケーション設定を表す構造体。
type Config struct {
	// Port ポート
	// 0 が指定された場合は動的にポートを割り当てる。
	Port uint16 `env:"PORT" envDefault:"8080"`

	// DBHost データベースホスト
	DBHost string `env:"DB_HOST" envDefault:"localhost"`
	// DBPort データベースポート
	DBPort uint16 `env:"DB_PORT" envDefault:"3306"`
	// DBName データベース名
	DBName string `env:"DB_NAME,required"`
	// DBUsername データベースユーザー名
	DBUsername string `env:"DB_USERNAME,required"`
	// DBPassword データベースパスワード
	DBPassword string `env:"DB_PASSWORD"`

	// AuthSecret 認証トークンの署名用シークレット
	AuthSecret string `env:"AUTH_SECRET" envDefault:"secret"`

	// AdminToolPath 管理ツールのローカルファイルシステム上のパス
	AdminToolLocalPath string `env:"ADMIN_TOOL_LOCAL_PATH"`
	// DocsLocalPath API ドキュメントのローカルファイルシステム上のパス
	DocsLocalPath string `env:"DOCS_LOCAL_PATH"`

	// ClientOrigin クライアントのオリジン
	ClientOrigin string `env:"CLIENT_ORIGIN"`

	// DebugCORS CORS デバッグモード
	DebugCORS bool `env:"DEBUG_CORS"`
	// FakeTime 時刻偽装モード設定
	// 時刻が指定された場合はその時刻に固定する。
	// Truthy な値が指定された場合はデフォルトの時刻に固定する。
	FakeTime FakeTimeMode `env:"FAKE_TIME"`
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
