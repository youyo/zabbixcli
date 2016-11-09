package main

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	latest "github.com/tcnksm/go-latest"
)

var Name string
var Version string

func main() {
	cli.VersionPrinter = func(c *cli.Context) { versionCheck() }
	app := cli.NewApp()
	app.Name = Name
	app.Version = Version
	app.Author = "youyo"
	app.Email = ""
	app.Usage = "For control the host of zabbix."

	app.Flags = GlobalFlags
	app.Commands = Commands
	app.CommandNotFound = CommandNotFound

	app.Run(os.Args)
}

func versionCheck() {
	githubTag := &latest.GithubTag{
		Owner:      "youyo",
		Repository: "zabbixcli",
	}
	res, err := latest.Check(githubTag, Version)
	if err == nil {
		if res.Outdated {
			fmt.Printf("%s is not latest, you should upgrade to %s\n", Version, res.Current)
		}
	} else {
		fmt.Printf("Network is not unreachable. Can not check version.\n")
	}
	fmt.Printf("%s version %s\n", Name, Version)
}
