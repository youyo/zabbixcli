package command

import (
	"errors"
	"log"

	"github.com/codegangsta/cli"
)

func CmdDelete(c *cli.Context) (err error) {
	// set logger
	setLoggerColog(c.GlobalBool("debug"))

	// Required args check
	m := map[string]string{
		"hostname": c.String("hostname"),
	}
	checkRequiredStringArgs(m)

	z := newZabbixctl(c)
	if err := z.login(); err != nil {
		log.Printf("error: %v", err)
		return err
	}

	exist, err := z.hostExists(c.String("hostname"))
	if err != nil {
		log.Printf("error: %v", err)
		return
	} else if exist == false {
		err = errors.New("Host is not exist.")
		log.Printf("error: %v", err)
		return
	}

	hostID, err := z.hostIdGet(c.String("hostname"))
	if err != nil {
		log.Printf("error: %v", err)
		return
	}

	err = z.hostDelete(hostID)
	if err != nil {
		log.Printf("error: %v", err)
		return
	} else {
		log.Printf("info: Host is deleted.")
	}

	return
}

func (z *zabbixcli) hostDelete(hostID string) (err error) {
	_, err = z.Api.Call("host.delete", []string{
		hostID,
	})
	return
}
