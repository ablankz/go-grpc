package config_test

import (
	"go-grpc/internal/config"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestGet(t *testing.T) {
	// Unset environment variables for test
	envKeys := []string{
		"PORT",
		"GO_ENV",
		"LOCAL_ROOT_PATH",
		"DEBUG",
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
				"LOCAL_ROOT_PATH": `/home/test/go-grpc`,
			},
			out: &config.Config{
				Port:      50051,
				Env:       "development",
				LocalRoot: "/home/test/go-grpc",
				Debug:     false,
			},
		},
		{
			name: "full",
			env: map[string]string{
				"PORT":            "50052",
				"GO_ENV":          "staging",
				"LOCAL_ROOT_PATH": "/home/test/go-grpc",
				"DEBUG":           "true",
			},
			out: &config.Config{
				Port:      50052,
				Env:       "staging",
				LocalRoot: "/home/test/go-grpc",
				Debug:     true,
			},
		},
		{
			name: "missing LOCAL_ROOT_PATH",
			env: map[string]string{
				"GO_ENV": "production",
				"DEBUG":  "true",
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
