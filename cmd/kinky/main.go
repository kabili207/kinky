package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"z0ne.dev/kura/kinky/sources/booru"

	"z0ne.dev/kura/kinky/sources/fs"

	"z0ne.dev/kura/kinky/config"

	"cdr.dev/slog"

	"cdr.dev/slog/sloggers/sloghuman"
	"github.com/urfave/cli/v2"
)

func before(ctx *cli.Context) error {
	configPath := ctx.Path("config")

	slogger := ctx.App.Metadata["log"].(slog.Logger)

	if ctx.Bool("verbose") {
		slogger = slogger.Leveled(slog.LevelDebug)
		ctx.App.Metadata["log"] = slogger
	}

	slogger.Info(context.Background(), "Hello!", slog.F("version", config.Version))
	slogger.Debug(context.Background(), "Git info", slog.F("branch", config.GitBranch), slog.F("hash", config.GitShaShort()), slog.F("dirty", config.GitDirty))

	_, err := os.Stat(configPath)
	if err == nil {
		cfg := new(config.Config)
		if err := cfg.Load(configPath); err != nil {
			return nil
		}

		ctx.App.Metadata["config"] = cfg
	}

	return nil
}

func main() {
	slogger := sloghuman.Make(os.Stdout)
	stdLog := slog.Stdlib(context.Background(), slogger)
	log.SetOutput(stdLog.Writer())

	fs.Register()
	booru.Register()

	app := &cli.App{
		Usage:       "kinky image bot",
		Description: "kinky is an image bot for the Mastodon or Pleroma network to post images",
		Compiled:    config.BuildTime(),
		Version:     config.Version,

		Copyright: "© 2020 z0ne",
		Authors: []*cli.Author{
			{Name: "Kura", Email: "me@kurabloodlust.eu"},
		},

		Metadata: map[string]interface{}{
			"log": slogger,
		},

		Flags: []cli.Flag{
			&cli.PathFlag{
				Name:      "config",
				Aliases:   []string{"c"},
				Usage:     "Specify config file to use",
				Value:     fmt.Sprintf("%v.yml", config.ApplicationName),
				TakesFile: true,
			},

			&cli.BoolFlag{
				Name:    "verbose",
				Aliases: []string{"V"},
				Usage:   "enable verbose output",
			},
		},

		Before: before,

		Commands: []*cli.Command{
			actionPost,
			actionInit,
			actionService,
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
