package backup

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"os"

	"github.com/cespare/xxhash/v2"
)

const BAKINFOILE = "bakinfo.json"

type BakInfo struct {
	Mtime int64  `json:"mtime"`
	Size  int64  `json:"size"`
	Hash  uint64 `json:"hash"`
}

type Bakinfo struct {
	Files   map[string]BakInfo `json:"files"`
}

var bakinfoContent Bakinfo

func init() {
	bakFile, err := os.Open(BAKINFOILE)
	if err != nil {
		switch {
		// create backinfo if not exist
		case errors.Is(err, os.ErrNotExist):
			{
				bakFile, err = os.Create(BAKINFOILE)
				bakFile.WriteString(`{"files": {}}`)
				if err != nil {
					slog.Error("backup init", "error", err)
				}
				slog.Info("backup init", "bakinfo", "created")
			}
		}
	}
	defer bakFile.Close()

	bakFileByte, err := io.ReadAll(bakFile)
	if err != nil {
		slog.Error("backup init", "error", err)
	}

	err = json.Unmarshal(bakFileByte, &bakinfoContent)
	if err != nil {
		slog.Error("backup init", "error", err)
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
	return os.WriteFile(BAKINFOILE, data, 0o644)
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
