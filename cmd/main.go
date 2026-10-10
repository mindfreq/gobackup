package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/mindfreq/gobackup/internal/backup"

	"github.com/urfave/cli/v3"
)

func main() {
	fmt.Println(backup.BAKINFOILE)
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
					_, err := os.Open(filePath)
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
