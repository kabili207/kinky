package main

import (
	"context"

	"z0ne.dev/kura/kinky/sources/fs"

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
			return fs.ErrNoConfig
		}

		log.Debug(context.Background(), "creating backend engine")
		s, err := cfg.ParseSource()
		if err != nil {
			return err
		}

		log.Debug(context.Background(), "fetching new image url")
		img, fileName, err := s.GetImageReader()
		if err != nil {
			return ErrNoImageFound
		}
		defer img.Close()

		log.Debug(context.Background(), "logging in to mastoderp")
		c := mastodon.NewClient(&mastodon.Config{
			Server:       cfg.Instance,
			ClientID:     cfg.ClientID,
			ClientSecret: cfg.ClientSecret,
			AccessToken:  cfg.AccessToken,
		})

		log.Info(context.Background(), "uploading image...", slog.F("file", img))
		att, err := c.UploadMediaFromReaderWithCustomFileName(context.Background(), img, fileName)
		if err != nil {
			return err
		}

		log.Debug(context.Background(), "generating new status")
		toot := &mastodon.Toot{
			Visibility: cfg.PostOptions.Visibility,
			Sensitive:  cfg.PostOptions.NSFW || s.IsSensitive(),
			Status:     cfg.PostOptions.Content,
			MediaIDs:   []mastodon.ID{att.ID},
		}

		msg, err := s.Caption()
		if err != nil {
			return err
		}

		if msg != "" {
			if cfg.PostOptions.AppendPostContent {
				toot.Status += "\n\n" + msg
			} else {
				toot.Status = msg
			}
		}

		log.Info(context.Background(), "posting new status...")
		st, err := c.PostStatus(context.Background(), toot)
		if err != nil {
			return err
		}

		log.Info(context.Background(), "status posted", slog.F("url", st.URL))

		return nil
	},
}
