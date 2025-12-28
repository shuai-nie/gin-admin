package cmd

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/urfave/cli/v2"
	"internal/bootstrap"
	"internal/config"
)

func StartCmd() *cli.Command {
	return &cli.Command{
		Name:  "start",
		Usage: "Start server",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "workdir",
				Aliases:     []string{"d"},
				Usage:       "Working directory",
				DefaultText: "configs",
				Value:       "configs",
			},
			&cli.StringFlag{
				Name:        "config",
				Aliases:     []string{"c"},
				Usage:       "Runtime configuration files or directory (relative to workdir, multiple separated by commas)",
				DefaultText: "dev",
				Value:       "dev",
			},
			&cli.StringFlag{
				Name:    "static",
				Aliases: []string{"s"},
				Usage:   "Static files directory",
			},
			&cli.BoolFlag{
				Name:  "daemon",
				Usage: "Run as a daemon",
			},
		},
		Action: func(c *cli.Context) error {
			workDir := c.String("workdir")
			staticDir := c.String("static")
			configs := c.String("config")

			if c.Bool("daemon") {
				bin, err := filepath.Abs(os.Args[0])
				if err != nil {
					fmt.Printf("faild to get absolute path for command: %s \n", err.Error())
					return err
				}

				args := []string{"start"}
				args = append(args, "-d", workDir)
				args = append(args, "-c", configs)
				args = append(args, "-s", staticDir)
				fmt.Printf("execute command:%s %s \n", bin, strings.Join(args, ""))
				command := exec.Command(bin, args...)

				stdLogFile := fmt.Sprintf("%s.log", c.App.Name)
				file, err := os.OpenFile(stdLogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
				if err != nil {
					fmt.Printf("failed to open log file: %s \n", err.Error())
					return err
				}
				defer file.Close()

				command.Stdout = file
				command.Stderr = file

				fmt.Printf("service %s daemon thread started started successfully\n", config.C.General.AppName)

				pid := command.Process.Pid
				_ = os.WriteFile(fmt.Sprintf("%s.lock", c.App.Name), []byte(fmt.Sprintf("%d", pid)), 066)
				fmt.Printf("service %s daemon thread started with %d \n", config.C.General.AppName, pid)
				os.Exit(0)

			}

			err := bootstrap.Run(context.Background(), bootstrap.RunConfig{
				WorkDir:   workDir,
				Configs:   configs,
				StaticDir: staticDir,
			})

			if err != nil {
				panic(err)
			}
			return nil
		},
	}
}
