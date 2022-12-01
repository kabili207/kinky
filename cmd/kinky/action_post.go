package main

import (
	"cdr.dev/slog"
	"github.com/urfave/cli/v2"
	"z0ne.dev/kura/kinky/config"
	"z0ne.dev/kura/kinky/sources/fs"
)

var actionPost = &cli.Command{
	Name:  "post",
	Usage: "Post a new message",

	Flags: nil,

	Action: func(ctx *cli.Context) error {
		ilog := ctx.App.Metadata["log"]
		log := ilog.(slog.Logger)

		icfg := ctx.App.Metadata["config"]
		cfg, ok := icfg.(*config.Config)
		if icfg == nil || !ok || cfg == nil {
			return fs.ErrNoConfig
		}
		return postImage(ctx, cfg, log)
	},
}
