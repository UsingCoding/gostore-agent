package main

import (
	"context"
	stdlog "log"
	"os"
	"os/signal"
	"path"
	"syscall"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

const (
	appID = "gostore-agent"
)

var (
	version = "dev"
	commit  = "none"
)

func main() {
	ctx := context.Background()

	ctx = subscribeForKillSignals(ctx)

	err := runApp(ctx, os.Args)
	if err != nil {
		stdlog.Fatal(err)
	}
}

func runApp(ctx context.Context, args []string) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = ""
	}

	a := &cli.App{
		Name:    appID,
		Version: version,
		// do not use built-in version flag
		HideVersion:          true,
		Usage:                "Agent plugin for gostore",
		EnableBashCompletion: true,
		Commands: []*cli.Command{
			versionCMD(),
			ssh(),
			install(),
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Value:   path.Join(homeDir, ".gostore-agent/config.toml"),
			},
			&cli.BoolFlag{
				Name:    "debug",
				Aliases: []string{"d"},
				Value:   false,
			},
		},
	}

	return a.RunContext(ctx, args)
}

func subscribeForKillSignals(ctx context.Context) context.Context {
	ctx, cancel := context.WithCancel(ctx)

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		defer cancel()
		select {
		case <-ctx.Done():
			signal.Stop(ch)
		case <-ch:
		}
	}()

	return ctx
}

func initLogger(ctx *cli.Context) *logrus.Logger {
	logger := logrus.New()
	if ctx.Bool("debug") {
		logger.SetLevel(logrus.DebugLevel)
	}
	logger.Debugln("KEK")
	return logger
}
