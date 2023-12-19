package config

import (
	"errors"
	"strconv"
	"time"
)

// DefaultFakeTime 時刻偽装モードのデフォルトの初期時刻。
var DefaultFakeTime = time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)

// FakeTimeMode 時刻偽装モードの設定を表す。
type FakeTimeMode struct {
	// Enabled true のとき有効
	Enabled bool
	// Time 初期時刻
	Time time.Time
}

// parseFakeTimeMode 時刻偽装モード設定文字列をパースする。
func parseFakeTimeMode(v string) (any, error) {
	if v == "" {
		return FakeTimeMode{}, nil
	}

	t, err := time.Parse(time.RFC3339, v)
	if err == nil {
		return FakeTimeMode{
			Enabled: true,
			Time:    t,
		}, nil
	}

	enabled, err := strconv.ParseBool(v)
	if err != nil {
		return nil, errors.New("invalid value")
	}

	if enabled {
		return FakeTimeMode{
			Enabled: true,
			Time:    DefaultFakeTime,
		}, nil
	}

	return FakeTimeMode{}, nil
}
