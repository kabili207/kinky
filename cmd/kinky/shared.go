package main

import (
	"context"
	"errors"
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

var (
	ErrAnimationsDisabled error = errors.New("animations are disabled")
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

	if strings.HasSuffix(fileName, ".mp4") {
		return ErrAnimationsDisabled
	}

	log.Debug(context.Background(), "logging in to misskey")

	mk, err := misskey.NewClientWithOptions(
		misskey.WithAPIToken(cfg.AccessToken),
		misskey.WithBaseURL("https", cfg.Instance, ""),
		misskey.WithLogLevel(logrus.InfoLevel),
	)
	if err != nil {
		return err
	}

	log.Info(context.Background(), "pulling image...", slog.F("file", fileName))

	fileContents, err := io.ReadAll(img)
	if err != nil {
		return err
	}

	hash := s.GetMd5Hash()
	matches, err := mk.Drive().File().FindByHash(hash)

	var file models.File

	if len(matches) > 0 {
		file = matches[0]
		log.Info(context.Background(), "image hash matched", slog.F("file", fileName), slog.F("file_id", file.ID))
	} else {
		log.Info(context.Background(), "uploading image...", slog.F("file", fileName))
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

	var cw core.String = nil
	switch rating := s.GetRating(); rating {
	case config.Sensitive:
		cw = core.NewString("Suggestive content")
	case config.Questionable:
		cw = core.NewString("Nudity and/or sexually suggestive content")
	case config.Explicit:
		cw = core.NewString("Sexually explicit content")
	}

	note := notes.CreateRequest{
		Text:       core.NewString(msg),
		Visibility: models.VisibilityPublic,
		FileIDs:    []string{file.ID},
		CW:         cw,
	}

	log.Info(context.Background(), "posting new status...")
	response, err := mk.Notes().Create(note)
	if err != nil {
		return err
	}

	log.Info(context.Background(), "status posted", slog.F("noteId", response.CreatedNote.ID))

	return nil
}
