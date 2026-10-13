package backup

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	// "path/filepath"

	"github.com/mindfreq/gobackup/internal/backup/services"
)

const BAKINFOILE = "bakinfo.json"

type FileInfo struct {
	Mtime int64  `json:"mtime"`
	Size  int64  `json:"size"`
	Hash  string `json:"hash"`
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

func IsChange(file *os.File) error {
	// absPath, err := filepath.Abs(file.Name())
	// if err != nil {
	// 	return err
	// }
	fileInfo := bakinfoContent.Files["/hol2/22"]
	// handle not exist
	if fileInfo.Size == 0 {
		
	} else {
		
	}
	fmt.Println(fileInfo)
	return nil
}

func PreBackup(ctx context.Context, svc services.Service ,file *os.File) error {
	info, err := file.Stat()
	if err != nil {
		return err
	}

	if info.IsDir() {
		// todo
	} else {
		err = svc.Upload(ctx, "test", file)
		if err != nil {
			return err
		}
		slog.Info("successful upload", "file", file.Name())
	}

	return nil
}
