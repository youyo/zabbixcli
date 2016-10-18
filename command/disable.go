package command

import (
	"errors"
	"log"

	"github.com/codegangsta/cli"
)

func CmdDisable(c *cli.Context) (err error) {
	// set logger
	setLoggerColog(c.GlobalBool("debug"))

	// Required args check
	m := map[string]string{
		"hostname": c.String("hostname"),
	}
	checkRequiredStringArgs(m)

	hostname := c.String("hostname")
	z := newZabbixctl(c)
	if err = z.login(); err != nil {
		log.Printf("error: %v", err)
		return
	}

	exist, err := z.hostExists(hostname)
	if err != nil {
		log.Printf("error: %v", err)
		return
	} else if exist == false {
		err = errors.New("Host is not exist.")
		log.Printf("error: %v", err)
		return err
	}

	hostID, err := z.hostIdGet(hostname)
	if err != nil {
		log.Printf("error: %v", err)
		return
	}

	_, err = z.hostStatusUpdate(hostID, HostStatusDisable)
	if err == nil {
		log.Printf("info: Host is disabled.")
	} else {
		log.Printf("error: %v", err)
	}
	return
}
