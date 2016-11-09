package command

import (
	"errors"
	"log"

	"github.com/AlekSi/zabbix"
	"github.com/codegangsta/cli"
)

func CmdSearch(c *cli.Context) (err error) {
	// set logger
	setLoggerColog(c.GlobalBool("debug"))

	// Required args check
	m := map[string]string{
		"hostname": c.String("hostname"),
	}
	checkRequiredStringArgs(m)

	searchWord := c.String("hostname")
	z := newZabbixctl(c)
	if err = z.login(); err != nil {
		log.Printf("error: %v", err)
		return
	}

	hosts, err := z.searchHosts(searchWord)
	if err != nil {
		log.Printf("error: %v", err)
		return
	} else {
		outputTable(hosts, "Host")
	}

	return
}

func (z *zabbixcli) searchHosts(searchWord string) (hosts []string, err error) {
	resp, err := z.Api.Call("host.get", zabbix.Params{
		"output": "extend",
		"search": map[string]string{
			"host": searchWord,
		},
	})
	if err != nil {
		return hosts, err
	}
	rr := resp.Result.([]interface{})
	for _, v := range rr {
		r, ok := v.(map[string]interface{})
		if !ok {
			err = errors.New("assertion error")
			return hosts, err
		}
		hosts = append(hosts, r["name"].(string))
	}
	return hosts, nil
}
