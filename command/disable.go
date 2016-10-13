package command

import (
	"errors"
	"fmt"

	"github.com/codegangsta/cli"
)

func CmdDisable(c *cli.Context) (err error) {
	// Required args check
	m := map[string]string{
		"hostname": c.String("hostname"),
	}
	checkRequiredStringArgs(m)

	hostname := c.String("hostname")
	z := newZabbixctl(c)
	if err = z.login(); err != nil {
		return
	}

	exist, err := z.hostExists(hostname)
	if err != nil {
		return
	} else if exist == false {
		return errors.New("Host is not exist.")
	}

	hostID, err := z.hostIdGet(hostname)
	if err != nil {
		return
	}

	_, err = z.hostStatusUpdate(hostID, HostStatusDisable)
	if err == nil {
		fmt.Println("Host is disabled.")
	}
	return
}
