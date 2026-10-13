package config

import (
	"io"
	"log/slog"
	"os"

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

func init() {
	configFile, err := os.Open("config.toml")
	if err != nil {
		slog.Error("config init", "error", err)
		os.Exit(1)
	}
	configFileByte, err := io.ReadAll(configFile)
	if err != nil {
		slog.Error("config init", "error", err)
		os.Exit(1)
	}

	err = toml.Unmarshal(configFileByte, &Config)
	if err != nil {
		slog.Error("config init", "error", err)
		os.Exit(1)
	}
}
