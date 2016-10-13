package command

import (
	"fmt"
	"os"
	"strings"

	"github.com/AlekSi/zabbix"
	"github.com/codegangsta/cli"
)

const HostStatusEnable = 0
const HostStatusDisable = 1
const UseIp = 1
const UseDns = 0

type (
	zabbixctl struct {
		Login login
		Api   *zabbix.API
	}

	login struct {
		UserName string
		Password string
		Url      string
	}
)

func ifEmptyString(s string) (b bool) {
	switch s {
	case "":
		b = true
		return
	default:
		b = false
		return
	}
}

func checkRequiredStringArgs(i map[string]string) {
	for k, v := range i {
		if ifEmptyString(v) {
			fmt.Fprintf(os.Stderr, "Arg: '%s' is must be required. See '--help'.", k)
			os.Exit(2)
		}
	}
}

func newZabbixctl(c *cli.Context) (z *zabbixctl) {
	// Required args check
	m := map[string]string{
		"username": c.GlobalString("username"),
		"password": c.GlobalString("password"),
		"url":      c.GlobalString("url"),
	}
	checkRequiredStringArgs(m)

	// Login
	z = &zabbixctl{
		Login: login{
			UserName: c.GlobalString("username"),
			Password: c.GlobalString("password"),
			Url:      c.GlobalString("url"),
		},
	}
	return
}

func (z *zabbixctl) login() (err error) {
	z.Api = zabbix.NewAPI(z.Login.Url)
	_, err = z.Api.Login(z.Login.UserName, z.Login.Password)
	return
}

func (z *zabbixctl) hostExists(host string) (exist bool, err error) {
	resp, err := z.Api.Call("host.exists", zabbix.Params{
		"host": host,
	})
	exist = resp.Result.(bool)
	return
}

func (z *zabbixctl) hostStatusUpdate(hostID string, status int) (resp zabbix.Response, err error) {
	resp, err = z.Api.Call("host.update", zabbix.Params{
		"hostid": hostID,
		"status": status,
	})
	return
}

func (z *zabbixctl) hostIdGet(host string) (hostID string, err error) {
	resp, err := z.Api.Call("host.get", zabbix.Params{
		"output": "extend",
		"filter": map[string][]string{
			"host": []string{
				host,
			},
		},
	})
	if err != nil {
		return "", err
	}
	hostID, _ = extractHostID(resp)
	return
}

func extractHostID(resp zabbix.Response) (hostID string, err error) {
	rr := resp.Result.([]interface{})
	r := rr[0].(map[string]interface{})
	hostID = r["hostid"].(string)
	err = nil
	return
}

func (z *zabbixctl) hostGroupIdGet(hostGroup string) (hostGroupID string, err error) {
	resp, err := z.Api.Call("hostgroup.get", zabbix.Params{
		"output": "extend",
		"filter": map[string][]string{
			"name": []string{
				hostGroup,
			},
		},
	})
	if err != nil {
		return "", err
	}
	hostGroupID, _ = extractHostGroupID(resp)
	return
}

func extractHostGroupID(resp zabbix.Response) (hostGroupID string, err error) {
	rr := resp.Result.([]interface{})
	r := rr[0].(map[string]interface{})
	hostGroupID = r["groupid"].(string)
	err = nil
	return
}

func (z *zabbixctl) templateIdGet(template string) (templateID string, err error) {
	resp, err := z.Api.Call("template.get", zabbix.Params{
		"output": "extend",
		"filter": map[string][]string{
			"host": []string{
				template,
			},
		},
	})
	if err != nil {
		return "", err
	}
	templateID, _ = extractTemplateID(resp)
	return
}

func extractTemplateID(resp zabbix.Response) (templateID string, err error) {
	rr := resp.Result.([]interface{})
	r := rr[0].(map[string]interface{})
	templateID = r["templateid"].(string)
	err = nil
	return
}

func splitTemplates(templates string) (t []string) {
	t = strings.Split(templates, ",")
	return
}

func (z *zabbixctl) proxyIdGet(proxy string) (proxyID string, err error) {
	resp, err := z.Api.Call("proxy.get", zabbix.Params{
		"output": "extend",
		"filter": map[string][]string{
			"host": []string{
				proxy,
			},
		},
	})
	if err != nil {
		return "", err
	}
	proxyID, _ = extractProxyID(resp)
	return
}

func extractProxyID(resp zabbix.Response) (proxyID string, err error) {
	rr := resp.Result.([]interface{})
	r := rr[0].(map[string]interface{})
	proxyID = r["proxyid"].(string)
	err = nil
	return
}
