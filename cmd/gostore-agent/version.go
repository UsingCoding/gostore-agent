package main

import (
	"encoding/json"
	"fmt"

	"github.com/urfave/cli/v2"
)

func versionCMD() *cli.Command {
	return &cli.Command{
		Name:   "version",
		Usage:  "Show gostore version",
		Action: executeVersion,
	}
}

func executeVersion(_ *cli.Context) error {
	v := struct {
		Version string `json:"version"`
		Commit  string `json:"commit"`
	}{
		Version: version,
		Commit:  commit,
	}
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}

	fmt.Println(string(data))

	return nil
}
