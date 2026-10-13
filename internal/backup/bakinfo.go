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

type FileInfo struct {
	Mtime int64  `json:"mtime"`
	Size  int64  `json:"size"`
	Hash  uint64 `json:"hash"`
}

type Bakinfo struct {
	Folders map[string]string   `json:"folders"`
	Files   map[string]FileInfo `json:"files"`
}

var bakinfoContent Bakinfo

func init() {
	backinfoFile, err := os.Open(BAKINFOILE)
	if err != nil {
		switch {
		// create backinfo if not exist
		case errors.Is(err, os.ErrNotExist):
			{
				backinfoFile, err = os.Create(BAKINFOILE)
				backinfoFile.WriteString(`{
	"folders": {},
	"files": {}
}`)
				if err != nil {
					slog.Error("backup init", "error", err)
				}
				slog.Info("backup init", "bakinfo", "created")
			}
		}
	}
	defer backinfoFile.Close()

	backinfoFileByte, err := io.ReadAll(backinfoFile)
	if err != nil {
		slog.Error("backup init", "error", err)
	}

	err = json.Unmarshal(backinfoFileByte, &bakinfoContent)
	if err != nil {
		slog.Error("backup init", "error", err)
	}

}

func hashFile(path string) (uint64, error) {
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
	info, err := os.Stat(path)
	if err != nil {
		return err
	}

	hash, err := hashFile(path)
	if err != nil {
		return err
	}

	bakinfoContent.Files[path] = FileInfo{
		Mtime: info.ModTime().UnixNano(),
		Size:  info.Size(),
		Hash:  hash,
	}
	return saveBakinfo()
}
