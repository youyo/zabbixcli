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
	cli.StringFlag{
		Name:   "username",
		Usage:  "set username.",
		EnvVar: "ZABBIXCTL_USERNAME",
	},
	cli.StringFlag{
		Name:   "password",
		Usage:  "set password.",
		EnvVar: "ZABBIXCTL_PASSWORD",
	},
	cli.StringFlag{
		Name:   "url",
		Usage:  "set url.",
		EnvVar: "ZABBIXCTL_URL",
	},
}

var Commands = []cli.Command{
	{
		Name:     "create",
		Category: "Host management",
		Usage:    "Create host",
		Action:   command.CmdCreate,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "hostname, H",
				Usage: "set hostname.",
			},
			cli.StringFlag{
				Name:  "host-group, g",
				Usage: "set host-group(s)",
			},
			cli.StringFlag{
				Name:  "templates, t",
				Usage: "set template(s).",
			},
			cli.StringFlag{
				Name:  "ipaddress, i",
				Usage: "set ipaddress",
			},
			cli.StringFlag{
				Name:  "dnsname, d",
				Usage: "set dnsname",
			},
			cli.BoolFlag{
				Name:  "use-ip",
				Usage: "select connect type. default false.(use dnsname)",
			},
			cli.StringFlag{
				Name:  "port, p",
				Usage: "set port",
				Value: "10050",
			},
			cli.StringFlag{
				Name:   "proxy, P",
				Usage:  "set proxy.",
				EnvVar: "ZABBIXCTL_PROXY",
			},
		},
	},
	{
		Name:     "enable",
		Category: "Host management",
		Usage:    "Enable host",
		Action:   command.CmdEnable,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "hostname, H",
				Usage: "set hostname.",
			},
		},
	},
	{
		Name:     "disable",
		Category: "Host management",
		Usage:    "Disable host",
		Action:   command.CmdDisable,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "hostname, H",
				Usage: "set hostname.",
			},
		},
	},
	{
		Name:     "delete",
		Category: "Host management",
		Usage:    "Delete host",
		Action:   command.CmdDelete,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "hostname, H",
				Usage: "set hostname.",
			},
		},
	},
	{
		Name:     "proxy-list",
		Category: "Lists",
		Usage:    "Shows a list of zabbix-proxies",
		Action:   command.CmdProxyList,
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "raw, r",
				Usage: "output raw format.",
			},
		},
	},
	{
		Name:     "template-list",
		Category: "Lists",
		Usage:    "Shows a list of templates",
		Action:   command.CmdTemplateList,
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "raw, r",
				Usage: "output raw format.",
			},
		},
	},
	{
		Name:     "hostgroup-list",
		Category: "Lists",
		Usage:    "Shows a list of host-groups",
		Action:   command.CmdHostGroupList,
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "raw, r",
				Usage: "output raw format.",
			},
		},
	},
	{
		Name:     "search",
		Category: "Lists",
		Usage:    "Search hosts",
		Action:   command.CmdSearch,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "hostname, H",
				Usage: "set hostname.",
			},
		},
	},
}

func CommandNotFound(c *cli.Context, command string) {
	fmt.Fprintf(os.Stderr, "%s: '%s' is not a %s command. See '%s --help'.", c.App.Name, command, c.App.Name, c.App.Name)
	os.Exit(2)
}
