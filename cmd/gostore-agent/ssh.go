package main

import (
	"github.com/urfave/cli/v2"

	"github.com/UsingCoding/gostore-agent/internal/service"
)

func ssh() *cli.Command {
	return &cli.Command{
		Name: "ssh",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "socket-path",
				Aliases:  []string{"s"},
				Usage:    "Path for ssh socket. ${HOME}/.gostore-agent/gostore-agent.sock",
				Required: true,
			},
		},
		Action: func(c *cli.Context) error {
			socketPath := c.String("socket-path")

			s := service.Service{
				Logger: initLogger(c),
			}

			return s.SSH(c.Context, service.SSHParams{
				ConfigPath: c.String("config"),
				SocketPath: socketPath,
			})
		},
	}
}
