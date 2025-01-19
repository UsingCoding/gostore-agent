package config

import (
	"io"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/pkg/errors"
)

type Config struct {
	Gostore string `toml:"gostore"`
	Stores  Stores `toml:"stores"`
}

func (c Config) GetGostore() string {
	s := os.ExpandEnv(c.Gostore)
	if s == "" {
		const defaultExec = "gostore"
		s = defaultExec
	}
	return s
}

type Stores map[string]Store

type Store struct {
	IDs []string `toml:"ids"`
}

func Parse(r io.ReadCloser) (c Config, err error) {
	defer func() {
		_ = r.Close()
	}()
	_, err = toml.NewDecoder(r).Decode(&c)

	return c, errors.Wrap(err, "failed to parse config")
}

func Dump(w io.Writer, c Config) error {
	return toml.NewEncoder(w).Encode(c)
}
