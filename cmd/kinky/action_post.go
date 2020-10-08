package main

import (
	"context"

	"z0ne.dev/kura/kinky/source"

	"cdr.dev/slog"
	"github.com/mattn/go-mastodon"

	"github.com/urfave/cli/v2"

	"z0ne.dev/kura/kinky/config"
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
			return source.ErrNoConfig
		}

		log.Debug(context.Background(), "building image fetcher")
		s, err := source.New(cfg)
		if err != nil {
			return err
		}

		log.Debug(context.Background(), "logging in to mastoderp")
		c := mastodon.NewClient(&mastodon.Config{
			Server:       cfg.Instance,
			ClientID:     cfg.ClientID,
			ClientSecret: cfg.ClientSecret,
			AccessToken:  cfg.AccessToken,
		})

		log.Debug(context.Background(), "generating new status")
		post, err := s.Post(func(file string) (mastodon.ID, error) {
			log.Info(context.Background(), "Uploading image...", slog.F("file", file))
			//nolint:govet
			att, err := c.UploadMedia(context.Background(), file)
			if err != nil {
				return mastodon.ID(0), err
			}

			return att.ID, nil
		})
		if err != nil {
			return err
		}

		log.Info(context.Background(), "Posting new status...")
		st, err := c.PostStatus(context.Background(), post)
		if err != nil {
			return err
		}

		log.Info(context.Background(), "Status posted", slog.F("url", st.URL))

		return err
	},
}
