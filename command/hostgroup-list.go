package command

import (
	"fmt"

	"github.com/AlekSi/zabbix"
	"github.com/codegangsta/cli"
)

func CmdHostGroupList(c *cli.Context) (err error) {
	z := newZabbixctl(c)
	if err = z.login(); err != nil {
		return
	}
	hostGroups, err := z.hostGroupGet()
	outputHostGroupList(hostGroups)
	return
}

func (z *zabbixctl) hostGroupGet() (hostGroups []string, err error) {
	resp, err := z.Api.Call("hostgroup.get", zabbix.Params{
		"output":    "extend",
		"sortfield": "name",
	})
	if err != nil {
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

func outputHostGroupList(hostGroups []string) {
	for _, v := range hostGroups {
		fmt.Println(v)
	}
}
