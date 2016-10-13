package command

import (
	"fmt"

	"github.com/AlekSi/zabbix"
	"github.com/codegangsta/cli"
)

func CmdProxyList(c *cli.Context) (err error) {
	z := newZabbixctl(c)
	if err = z.login(); err != nil {
		return
	}
	proxies, err := z.proxyGet()
	if err != nil {
		return
	}
	outputProxyList(proxies)
	return
}

func (z *zabbixctl) proxyGet() (proxies []string, err error) {
	resp, err := z.Api.Call("proxy.get", zabbix.Params{
		"output":    "extend",
		"sortfield": "host",
	})
	if err != nil {
		return nil, err
	}
	proxies, _ = extractProxyHost(resp)
	return
}

func extractProxyHost(resp zabbix.Response) (proxies []string, err error) {
	rr := resp.Result.([]interface{})
	for _, r := range rr {
		r := r.(map[string]interface{})
		proxies = append(proxies, r["host"].(string))
	}
	return proxies, nil
}

func outputProxyList(proxies []string) {
	for _, v := range proxies {
		fmt.Println(v)
	}
}
