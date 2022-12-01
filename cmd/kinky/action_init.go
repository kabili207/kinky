package main

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"

	"z0ne.dev/kura/kinky/sources/fs"

	"github.com/manifoldco/promptui"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"github.com/yitsushi/go-misskey"

	"z0ne.dev/kura/kinky/config"
)

func validateURL(input string) error {
	_, err := url.Parse(input)

	return err
}

func ask(question string, validate promptui.ValidateFunc) string {
	prompt := promptui.Prompt{Label: question, Validate: validate}
	var err error
	var answer string
	for len(answer) == 0 {
		answer, err = prompt.Run()
		if err != nil {
			if errors.Is(err, promptui.ErrEOF) || errors.Is(err, promptui.ErrInterrupt) {
				os.Exit(0)
			}

			fmt.Printf("Invalid response: %v\nPlease try again...", err)
		}
	}

	return answer
}

func sel(question string, answers []string) int {
	prompt := promptui.Select{
		Label: question,
		Items: answers,
	}
	var err error
	answer := -1
	for answer < 0 {
		answer, _, err = prompt.Run()
		if err != nil {
			log.Printf("Invalid response: %v\nPlease try again...", err)
			answer = -1
		}
	}

	return answer
}

var actionInit = &cli.Command{
	Category: "Setup",
	Name:     "init",
	Usage:    "Initialize the bot",

	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "instance",
			Aliases: []string{"i"},
			Usage:   "URL of instance",
		},

		&cli.PathFlag{
			Name:    "folder",
			Aliases: []string{"fs", "fs-folder", "sf", "f"},
			Usage:   "Folder with images",
		},
	},

	Action: func(ctx *cli.Context) error {
		cfg := new(config.Config)

		cfg.Instance = ctx.String("instance")
		if len(cfg.Instance) == 0 {
			cfg.Instance = ask("Please enter the instance address", validateURL)
		}
		if len(cfg.AccessToken) == 0 {
			cfg.AccessToken = ask("Please enter the API Key", nil)
		}

		_, err := misskey.NewClientWithOptions(
			misskey.WithAPIToken(cfg.AccessToken),
			misskey.WithBaseURL("https", cfg.Instance, ""),
			misskey.WithLogLevel(logrus.DebugLevel),
		)

		if err != nil {
			return err
		}

		//fmt.Printf("Authenticated as %v\n", mk.)

		cfg.PostOptions.Visibility = "unlisted"
		cfg.PostOptions.NSFW = false
		cfg.PostOptions.Content = "."
		cfg.PostOptions.AppendPostContent = true

		sourceCfg := new(fs.SourceConfig)

		sourceCfg.Folder = ctx.String("folder")
		if len(sourceCfg.Folder) == 0 {
			sourceCfg.Folder = ask("Folder with images", nil)
		}
		sourceCfg.Recursive = true
		sourceCfg.EnableNSFWSuffix = true
		sourceCfg.EnableContentText = true
		sourceCfg.Extensions = []string{
			"png",
			"jpg",
			"jpeg",
			"gif",
			"webp",
			"mp4",
			"webm",
		}

		if err = cfg.SaveTo(ctx.Path("config")); err != nil {
			return err
		}

		fmt.Println("You are good to go!")

		return nil
	},
}
