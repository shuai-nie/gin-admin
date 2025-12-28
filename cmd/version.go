package cmd

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

func VersionCmd(v string) *cli.Command {
	return &cli.Command{
		Name:  "version",
		Usage: "show version",
		Action: func(_ *cli.Context) error {
			fmt.Printf(v)
			return nil
		},
	}
}
