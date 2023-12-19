package config_test

import (
	"go-grpc/internal/config"
	"os"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestGet(t *testing.T) {
	// Unset environment variables for test
	envKeys := []string{
		"PORT",
		"DB_HOST",
		"DB_PORT",
		"DB_NAME",
		"DB_USERNAME",
		"DB_PASSWORD",
		"AUTH_SECRET",
		"ADMIN_TOOL_LOCAL_PATH",
		"DOCS_LOCAL_PATH",
		"CLIENT_ORIGIN",
		"DEBUG_CORS",
		"FAKE_TIME",
	}
	for _, v := range envKeys {
		t.Setenv(v, "")
		os.Unsetenv(v)
	}

	cases := []struct {
		name   string
		env    map[string]string
		out    *config.Config
		failed bool
	}{
		{
			name: "minimum",
			env: map[string]string{
				"DB_NAME":     "app",
				"DB_USERNAME": "user",
			},
			out: &config.Config{
				Port:       8080,
				DBHost:     "localhost",
				DBPort:     3306,
				DBName:     "app",
				DBUsername: "user",
				AuthSecret: "secret",
			},
		},
		{
			name: "full",
			env: map[string]string{
				"PORT":                  "3000",
				"DB_HOST":               "db",
				"DB_PORT":               "9999",
				"DB_NAME":               "app",
				"DB_USERNAME":           "user",
				"DB_PASSWORD":           "password",
				"AUTH_SECRET":           "mysecret",
				"ADMIN_TOOL_LOCAL_PATH": "admin/dist",
				"DOCS_LOCAL_PATH":       "openapi",
				"CLIENT_ORIGIN":         "http://localhost:12345",
				"DEBUG_CORS":            "true",
				"FAKE_TIME":             "true",
			},
			out: &config.Config{
				Port:               3000,
				DBHost:             "db",
				DBPort:             9999,
				DBName:             "app",
				DBUsername:         "user",
				DBPassword:         "password",
				AuthSecret:         "mysecret",
				AdminToolLocalPath: "admin/dist",
				DocsLocalPath:      "openapi",
				ClientOrigin:       "http://localhost:12345",
				DebugCORS:          true,
				FakeTime: config.FakeTimeMode{
					Enabled: true,
					Time:    config.DefaultFakeTime,
				},
			},
		},
		{
			name: "FAKE_TIME is RFC3339 string",
			env: map[string]string{
				"DB_NAME":     "app",
				"DB_USERNAME": "user",
				"FAKE_TIME":   "2023-01-02T12:34:56Z",
			},
			out: &config.Config{
				Port:       8080,
				DBHost:     "localhost",
				DBPort:     3306,
				DBName:     "app",
				DBUsername: "user",
				AuthSecret: "secret",
				FakeTime: config.FakeTimeMode{
					Enabled: true,
					Time:    time.Date(2023, 1, 2, 12, 34, 56, 0, time.UTC),
				},
			},
		},
		{
			name: "FAKE_TIME is true",
			env: map[string]string{
				"DB_NAME":     "app",
				"DB_USERNAME": "user",
				"FAKE_TIME":   "true",
			},
			out: &config.Config{
				Port:       8080,
				DBHost:     "localhost",
				DBPort:     3306,
				DBName:     "app",
				DBUsername: "user",
				AuthSecret: "secret",
				FakeTime: config.FakeTimeMode{
					Enabled: true,
					Time:    config.DefaultFakeTime,
				},
			},
		},
		{
			name: "FAKE_TIME is 1",
			env: map[string]string{
				"DB_NAME":     "app",
				"DB_USERNAME": "user",
				"FAKE_TIME":   "1",
			},
			out: &config.Config{
				Port:       8080,
				DBHost:     "localhost",
				DBPort:     3306,
				DBName:     "app",
				DBUsername: "user",
				AuthSecret: "secret",
				FakeTime: config.FakeTimeMode{
					Enabled: true,
					Time:    config.DefaultFakeTime,
				},
			},
		},
		{
			name: "FAKE_TIME is false",
			env: map[string]string{
				"DB_NAME":     "app",
				"DB_USERNAME": "user",
				"FAKE_TIME":   "false",
			},
			out: &config.Config{
				Port:       8080,
				DBHost:     "localhost",
				DBPort:     3306,
				DBName:     "app",
				DBUsername: "user",
				AuthSecret: "secret",
			},
		},
		{
			name: "FAKE_TIME is 0",
			env: map[string]string{
				"DB_NAME":     "app",
				"DB_USERNAME": "user",
				"FAKE_TIME":   "0",
			},
			out: &config.Config{
				Port:       8080,
				DBHost:     "localhost",
				DBPort:     3306,
				DBName:     "app",
				DBUsername: "user",
				AuthSecret: "secret",
			},
		},
		{
			name: "FAKE_TIME is empty string",
			env: map[string]string{
				"DB_NAME":     "app",
				"DB_USERNAME": "user",
				"FAKE_TIME":   "",
			},
			out: &config.Config{
				Port:       8080,
				DBHost:     "localhost",
				DBPort:     3306,
				DBName:     "app",
				DBUsername: "user",
				AuthSecret: "secret",
			},
		},
		{
			name: "invalid PORT",
			env: map[string]string{
				"PORT":        "invalid",
				"DB_NAME":     "app",
				"DB_USERNAME": "user",
			},
			failed: true,
		},
		{
			name: "invalid DB_PORT",
			env: map[string]string{
				"DB_PORT":     "invalid",
				"DB_NAME":     "app",
				"DB_USERNAME": "user",
			},
			failed: true,
		},
		{
			name: "invalid FAKE_TIME",
			env: map[string]string{
				"DB_NAME":     "app",
				"DB_USERNAME": "user",
				"FAKE_TIME":   "invalid",
			},
			failed: true,
		},
		{
			name: "missing DB_NAME",
			env: map[string]string{
				"DB_USERNAME": "user",
			},
			failed: true,
		},
		{
			name: "missing DB_USERNAME",
			env: map[string]string{
				"DB_NAME": "app",
			},
			failed: true,
		},
	}

	for _, v := range cases {
		t.Run(v.name, func(tt *testing.T) {
			for key, value := range v.env {
				tt.Setenv(key, value)
			}

			cfg, err := config.Get()
			switch {
			case err != nil && !v.failed:
				tt.Fatalf("unexpected error: %+v", err)
			case err == nil && v.failed:
				tt.Fatal("unexpected success")
			case err != nil && v.failed:
				// pass
				tt.Logf("expected error: %+v", err)
				return
			}

			if diff := cmp.Diff(v.out, cfg); diff != "" {
				tt.Errorf("unexpected result:\n%s", diff)
			}
		})
	}
}
