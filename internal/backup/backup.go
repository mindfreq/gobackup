package backup

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"
	"sync"
	
	"golang.org/x/sync/errgroup"
	
	"github.com/mindfreq/gobackup/internal/backup/services"
)

func hasChange(path string, info os.FileInfo) (bool, error) {
	old, ok := bakinfoContent.Files[path]
	if !ok {
		// file not found (chenge it to add)
		return true, nil
	}

	if old.Mtime == info.ModTime().UnixNano() && old.Size == info.Size() {
		return false, nil
	}
	hash, err := GethashFile(path)
	if err != nil {
		return false, err
	}
	return hash != old.Hash, nil
}

var bakMu sync.Mutex

func handleFile(ctx context.Context, svc services.Service, file *os.File) error {
	info, err := os.Stat(file.Name())
	if err != nil {
		return err
	}
	absPath, err := filepath.Abs(file.Name())
	if err != nil {
		return err
	}

	// prevent race
	bakMu.Lock()
	chenged, err := hasChange(absPath, info)
	bakMu.Unlock()

	if err != nil {
		return err
	}
	if !chenged {
		slog.Info("file has not changed", "file", absPath)
		return nil
	}

	if err := svc.Upload(ctx, "backup"+absPath, file); err != nil {
		return err
	}
	if err := editBakfile(absPath); err != nil {
		return err
	}
	slog.Info("successful upload", "file", absPath)
	return nil
}

func handleFolder(ctx context.Context, svc services.Service, dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	var g errgroup.Group
	g.SetLimit(50)
	for _, e := range entries {
		g.Go(func() error {
			full := filepath.Join(dir, e.Name())
			if e.IsDir() {
				if err := handleFolder(ctx, svc, full); err != nil {
					return err
				}
			} else {
				f, err := os.Open(full)
				if err != nil {
					return err
				}
				err = handleFile(ctx, svc, f)
				f.Close()
				if err != nil {
					return err
				}
			}
			return nil
		})
	}
	if err := g.Wait(); err != nil {
		return err
	}
	return nil
}

func PreBackup(ctx context.Context, svc services.Service, file *os.File) error {
	info, err := file.Stat()
	if err != nil {
		return err
	}
	if info.IsDir() {
		path, err := filepath.Abs(file.Name())
		if err != nil {
			return err
		}
		handleFolder(ctx, svc, path)
		return nil
	} else {
		return handleFile(ctx, svc, file)
	}
}
