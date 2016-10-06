package main

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	"github.com/youyo/zabbixctl/command"
)

var GlobalFlags = []cli.Flag{
	cli.BoolFlag{
		Name:  "debug",
		Usage: "Set LogLevel Debug.",
	},
}

var Commands = []cli.Command{
	{
		Name:   "create",
		Usage:  "",
		Action: command.CmdCreate,
		Flags:  []cli.Flag{},
	},
	{
		Name:   "delete",
		Usage:  "",
		Action: command.CmdDelete,
		Flags:  []cli.Flag{},
	},
}

func CommandNotFound(c *cli.Context, command string) {
	fmt.Fprintf(os.Stderr, "%s: '%s' is not a %s command. See '%s --help'.", c.App.Name, command, c.App.Name, c.App.Name)
	os.Exit(2)
}