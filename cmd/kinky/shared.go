package main

import (
	"context"
	"html"
	"io"
	"strings"

	"cdr.dev/slog"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"github.com/yitsushi/go-misskey"
	"github.com/yitsushi/go-misskey/core"
	"github.com/yitsushi/go-misskey/models"
	"github.com/yitsushi/go-misskey/services/drive/files"
	"github.com/yitsushi/go-misskey/services/notes"
	"z0ne.dev/kura/kinky/config"
)

func postImage(ctx *cli.Context, cfg *config.Config, log slog.Logger) error {

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

	log.Debug(context.Background(), "logging in to misskey")

	mk, err := misskey.NewClientWithOptions(
		misskey.WithAPIToken(cfg.AccessToken),
		misskey.WithBaseURL("https", cfg.Instance, ""),
		misskey.WithLogLevel(logrus.InfoLevel),
	)
	if err != nil {
		return err
	}

	log.Info(context.Background(), "uploading image...", slog.F("file", fileName))

	fileContents, err := io.ReadAll(img)
	if err != nil {
		return err
	}

	// TODO: Check existance first
	hash := s.GetMd5Hash()
	matches, err := mk.Drive().File().FindByHash(hash)

	var file models.File

	if len(matches) > 0 {
		file = matches[0]
	} else {
		file, err = mk.Drive().File().Create(files.CreateRequest{
			FolderID:    cfg.PostOptions.FolderID,
			Name:        fileName,
			IsSensitive: cfg.PostOptions.NSFW || s.IsSensitive(),
			Force:       false,
			Content:     fileContents,
		})
	}
	if err != nil {
		return err
	}

	log.Debug(context.Background(), "generating new status")

	msg, err := s.Caption()
	if err != nil {
		return err
	}

	msg = html.UnescapeString(strings.TrimSpace(msg))

	if msg != "" {
		if cfg.PostOptions.AppendPostContent && cfg.PostOptions.Content != "" {
			msg = cfg.PostOptions.Content + "\n\n" + msg
		}
	} else {
		msg = cfg.PostOptions.Content
	}

	note := notes.CreateRequest{
		Text:       core.NewString(msg),
		Visibility: models.VisibilityPublic,
		FileIDs:    []string{file.ID},
	}

	log.Info(context.Background(), "posting new status...")
	response, err := mk.Notes().Create(note)
	if err != nil {
		return err
	}

	log.Info(context.Background(), "status posted", slog.F("noteId", response.CreatedNote.ID))

	return nil
}
