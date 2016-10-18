package command

import (
	"log"

	"github.com/AlekSi/zabbix"
	"github.com/codegangsta/cli"
)

func CmdHostGroupList(c *cli.Context) (err error) {
	// set logger
	setLoggerColog(c.GlobalBool("debug"))

	z := newZabbixctl(c)
	if err = z.login(); err != nil {
		log.Printf("error: %v", err)
		return
	}
	hostGroups, err := z.hostGroupGet()
	if err != nil {
		log.Printf("error: %v", err)
		return
	}

	// select output format
	switch c.Bool("raw") {
	case true:
		outputRaw(hostGroups)
	default:
		outputTable(hostGroups, "HostGroups")
	}
	return
}

func (z *zabbixctl) hostGroupGet() (hostGroups []string, err error) {
	resp, err := z.Api.Call("hostgroup.get", zabbix.Params{
		"output":    "extend",
		"sortfield": "name",
	})
	if err != nil {
		log.Printf("error: %v", err)
		return nil, err
	}
	hostGroups, _ = extractHostGroupName(resp)
	return
}

func extractHostGroupName(resp zabbix.Response) (hostGroups []string, err error) {
	rr := resp.Result.([]interface{})
	for _, r := range rr {
		r := r.(map[string]interface{})
		hostGroups = append(hostGroups, r["name"].(string))
	}
	return hostGroups, nil
}
