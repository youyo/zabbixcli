package command

import (
	"log"

	"github.com/AlekSi/zabbix"
	"github.com/codegangsta/cli"
)

func CmdProxyList(c *cli.Context) (err error) {
	// set logger
	setLoggerColog(c.GlobalBool("debug"))

	z := newZabbixctl(c)
	if err = z.login(); err != nil {
		log.Printf("error: %v", err)
		return
	}
	proxies, err := z.proxyGet()
	if err != nil {
		log.Printf("error: %v", err)
		return
	}

	// select output format
	switch c.Bool("raw") {
	case true:
		outputRaw(proxies)
	default:
		outputTable(proxies, "Proxies")
	}
	return
}

func (z *zabbixcli) proxyGet() (proxies []string, err error) {
	resp, err := z.Api.Call("proxy.get", zabbix.Params{
		"output":    "extend",
		"sortfield": "host",
	})
	if err != nil {
		log.Printf("error: %v", err)
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
