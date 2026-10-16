package config

import (
	"errors"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type s3 struct {
	Endpoint string
	Bucket   string
	Region   string
}

type config struct {
	S3 s3
}

var Config config

var DefaultConfigPath string

func init() {
	ConfigDir, err := os.UserConfigDir()
	if err != nil {
		homeDir, _ := os.UserHomeDir()
		DefaultConfigPath = filepath.Join(homeDir, "gobackup", "config.toml")
	} else {
		DefaultConfigPath = filepath.Join(ConfigDir, "gobackup", "config.toml")	
	}

	data, err := os.ReadFile(DefaultConfigPath)
	
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			slog.Error("config read", "error", err)
			os.Exit(1)
		}
		if err := os.MkdirAll(filepath.Dir(DefaultConfigPath), 0o755); err != nil {
	        slog.Error("config mkdir", "error", err)
	        os.Exit(1)
	    }
	    data = []byte("[s3]\n")
	    if err := os.WriteFile(DefaultConfigPath, data, 0o600); err != nil {
	        slog.Error("config init", "error", err)
	        os.Exit(1)
	    }
	}

	err = toml.Unmarshal(data, &Config)
	if err != nil {
		slog.Error("config init", "error", err)
		os.Exit(1)
	}
}
