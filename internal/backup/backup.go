package backup

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/mindfreq/gobackup/internal/backup/services"
)

func IsChange(file *os.File) (bool, error) {
	fileInfo, err := os.Stat(file.Name())
	if err != nil {
		return false, err
	}

	absPath, err := filepath.Abs(file.Name())
	if err != nil {
		return false, err
	}

	fileContent := bakinfoContent.Files[absPath]

	// handle not exist
	if fileContent.Size == 0 {
		err = editBakfile(absPath)
		if err != nil {
			return false, err
		}
		return true, nil
	} else { // check chnaged if file exist
		if fileContent.Mtime != fileInfo.ModTime().UnixNano() || fileContent.Size != fileInfo.Size() {
			// file was changed
			err = editBakfile(absPath)
			if err != nil {
				return true, err
			}
			return true, nil
		} else {
			// file not change
			return false, nil
		}
	}
}

func PreBackup(ctx context.Context, svc services.Service, file *os.File) error {
	info, err := file.Stat()
	if err != nil {
		return err
	}

	if info.IsDir() {
		// todo
	} else { // if a file
		isChange, err := IsChange(file)
		if err != nil {
			return err
		}
		abspath, err := filepath.Abs(file.Name())
		if err != nil {
			return err
		}
		if isChange {
			key := "backup" + abspath
			err = svc.Upload(ctx, key, file)
			if err != nil {
				return err
			}
			slog.Info("successful upload", "file", abspath)
		} else {
			// if file not change
			slog.Info("file was not changed", "file", abspath)
			return nil
		}
	}
	return nil
}
