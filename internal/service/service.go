package service

import (
	"context"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/UsingCoding/gostore-agent/internal/config"
	"github.com/UsingCoding/gostore-agent/internal/gostore"
	"github.com/UsingCoding/gostore-agent/internal/sshagent"
)

type Service struct {
	Logger logrus.FieldLogger
}

type SSHParams struct {
	ConfigPath string
	SocketPath string
}

func (s Service) SSH(ctx context.Context, params SSHParams) error {
	f, err := os.Open(params.ConfigPath)
	if err != nil {
		if os.IsNotExist(err) {
			return errors.Errorf("config file %s not found", params.ConfigPath)
		}

		return errors.Wrapf(err, "opening config file at %s", params.ConfigPath)
	}

	c, err := config.Parse(f)
	if err != nil {
		return errors.Wrapf(err, "parsing config file at %s", params.ConfigPath)
	}

	a := sshagent.NewAgent(
		c.Stores,
		gostore.NewLogAPI(
			gostore.NewGostore(c.GetGostore()),
			s.Logger,
		),
	)
	a = sshagent.NewLogAgent(a, s.Logger)

	return sshagent.Listen(
		ctx,
		os.ExpandEnv(params.SocketPath),
		a,
		s.Logger,
	)
}

type InstallParams struct {
	ConfigPath string
}

func (s Service) Install(params InstallParams) error {
	if exists(params.ConfigPath) {
		s.Logger.Infoln("Config already exists, skip")
		return nil
	}

	dir := filepath.Dir(params.ConfigPath)
	err := os.MkdirAll(dir, 0o755)
	if err != nil {
		return errors.Wrapf(err, "failed to create config dir %s", dir)
	}

	f, err := os.Create(params.ConfigPath)
	if err != nil {
		return errors.Wrapf(err, "failed to open config file %s", params.ConfigPath)
	}
	defer func() { _ = f.Close() }()

	err = config.Dump(
		f,
		config.Config{
			Stores: map[string]config.Store{
				"example": {},
			},
		},
	)
	return errors.Wrap(err, "failed to write config file")
}

func exists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	// Ensure it's not a directory
	return !info.IsDir()
}
