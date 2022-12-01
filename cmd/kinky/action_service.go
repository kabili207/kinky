package main

import (
	"context"
	"time"

	"z0ne.dev/kura/kinky/sources/fs"

	"cdr.dev/slog"

	"github.com/urfave/cli/v2"

	"z0ne.dev/kura/kinky/config"
)

var actionService = &cli.Command{
	Name:  "service",
	Usage: "Start the post service",

	Flags: nil,

	Action: runServiceLoop,
}

func runServiceLoop(ctx *cli.Context) error {
	ilog := ctx.App.Metadata["log"]
	log := ilog.(slog.Logger)

	icfg := ctx.App.Metadata["config"]
	cfg, ok := icfg.(*config.Config)
	if icfg == nil || !ok || cfg == nil {
		return fs.ErrNoConfig
	}

	delay := time.Duration(cfg.RunInterval) * time.Minute
	timer := time.NewTimer(delay)
	postImage(ctx, cfg, log)
	defer timer.Stop()
	for {
		timer.Reset(delay)
		select {
		case <-ctx.Done():
			log.Debug(context.Background(), "Breaking out of service loop:")
			return ctx.Err()
		case <-timer.C:
			postImage(ctx, cfg, log)
		}
	}
}
