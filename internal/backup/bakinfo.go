package backup

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/cespare/xxhash/v2"
)

var BakInfoFile string

type BakInfo struct {
	Mtime int64  `json:"mtime"`
	Size  int64  `json:"size"`
	Hash  uint64 `json:"hash"`
}

type Bakinfo struct {
	Files map[string]BakInfo `json:"files"`
}
var bakinfoContent Bakinfo

func init() {
	ConfigDir, err := os.UserConfigDir()
	if err != nil {
		homeDir, _ := os.UserHomeDir()
		BakInfoFile = filepath.Join(homeDir, "gobackup", "bakinfo.json")
	} else {
		BakInfoFile = filepath.Join(ConfigDir, "gobackup", "bakinfo.json")	
	}
	data, err := os.ReadFile(BakInfoFile)
	switch {
	// create backinfo if not exist
	case errors.Is(err, os.ErrNotExist):
		data = []byte(`{"files": {}}`)
		if err := os.WriteFile(BakInfoFile, data, 0o644); err != nil {
			slog.Error("backup init", "error", err)
			return
		}
		slog.Info("backup init", "bakinfo", "created")

	case err != nil:
		slog.Error("backup init", "error", err)
		return
	}

	err = json.Unmarshal(data, &bakinfoContent)
	if err != nil {
		slog.Error("backup init", "error", err)
		os.Exit(1)
	}
}

func GethashFile(path string) (uint64, error) {
	f, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	h := xxhash.New()
	if _, err := io.Copy(h, f); err != nil {
		return 0, err
	}
	return h.Sum64(), nil
}

func saveBakinfo() error {
	data, err := json.MarshalIndent(bakinfoContent, "", "\t")
	if err != nil {
		return err
	}
	return os.WriteFile(BakInfoFile, data, 0o644)
}

func editBakfile(path string) error {
	bakMu.Lock()
	defer bakMu.Unlock()
	info, err := os.Stat(path)
	if err != nil {
		return err
	}

	hash, err := GethashFile(path)
	if err != nil {
		return err
	}

	bakinfoContent.Files[path] = BakInfo{
		Mtime: info.ModTime().UnixNano(),
		Size:  info.Size(),
		Hash:  hash,
	}
	return saveBakinfo()
}
