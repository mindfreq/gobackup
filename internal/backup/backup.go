package backup

import (
	"encoding/json"
	"errors"


	"io"
	"log/slog"
	"os"
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

	var bakinfoContent Bakinfo
	err = json.Unmarshal(backinfoFileByte, &bakinfoContent)
	if err != nil {
		slog.Error("backup init", "error", err)
	}

}

// func AddFD(filePath os.File) {
// 	filePath
// }
