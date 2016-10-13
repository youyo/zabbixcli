package main

import (
	"os"

	"github.com/codegangsta/cli"
)

var Name string
var Version string

func main() {

	app := cli.NewApp()
	app.Name = Name
	app.Version = Version
	app.Author = "youyo"
	app.Email = ""
	app.Usage = "To controle the host of zabbix."

	app.Flags = GlobalFlags
	app.Commands = Commands
	app.CommandNotFound = CommandNotFound

	app.Run(os.Args)
}
