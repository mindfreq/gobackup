package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/mindfreq/gobackup/internal/backup"
	"github.com/mindfreq/gobackup/internal/backup/services"
	"github.com/urfave/cli/v3"
)

func main() {
	ctx := context.Background()
	svc, err := services.NewS3(ctx)
	if err != nil {
		slog.Error("s3 load", "error", err)
		os.Exit(1)
	}

	cmd := &cli.Command{
		Name:  "gobackup",
		Usage: "Backup Tool",
		Flags: []cli.Flag{},
		Commands: []*cli.Command{
			{
				Name:    "add",
				Aliases: []string{"a"},
				Usage:   "Add file or folder",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					filePath := cmd.Args().First()
					f, err := os.Open(filePath)
					if err != nil {
						slog.Error("main", "error", err)
						return nil
					}
					defer f.Close()
					backup.IsChange(f)

					err = backup.PreBackup(ctx, svc ,f)
					if err != nil {
						slog.Error("main", "error", err)
						return nil
					}

					return nil
				},
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		slog.Error("main", "error", err)
	}
}
