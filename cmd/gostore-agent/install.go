package main

import (
	"github.com/urfave/cli/v2"

	"github.com/UsingCoding/gostore-agent/internal/service"
)

func install() *cli.Command {
	return &cli.Command{
		Name: "install",
		Action: func(c *cli.Context) error {
			s := service.Service{
				Logger: initLogger(c),
			}

			return s.Install(service.InstallParams{
				ConfigPath: c.String("config"),
			})
		},
	}
}
