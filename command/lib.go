package command

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/AlekSi/zabbix"
	"github.com/codegangsta/cli"
	"github.com/comail/colog"
	"github.com/olekukonko/tablewriter"
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
	resp, err := z.Api.Call("host.get", zabbix.Params{
		"output": "extend",
		"filter": map[string][]string{
			"host": []string{
				host,
			},
		},
	})
	rr := resp.Result.([]interface{})
	size := len(rr)
	switch size {
	case 0:
		exist = false
		return
	default:
		exist = true
		return
	}
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
	hostID, err = extractHostID(resp)
	if err != nil {
		return "", err
	}
	return
}

func extractHostID(resp zabbix.Response) (hostID string, err error) {
	r, err := assertFirstResult(resp)
	if err != nil {
		return "", err
	}
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
	hostGroupID, err = extractHostGroupID(resp)
	if err != nil {
		return "", err
	}
	return
}

func extractHostGroupID(resp zabbix.Response) (hostGroupID string, err error) {
	r, err := assertFirstResult(resp)
	if err != nil {
		return "", err
	}
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
	templateID, err = extractTemplateID(resp)
	if err != nil {
		return "", err
	}
	return
}

func extractTemplateID(resp zabbix.Response) (templateID string, err error) {
	r, err := assertFirstResult(resp)
	if err != nil {
		return "", err
	}
	templateID = r["templateid"].(string)
	err = nil
	return
}

func splitArgs(args string) (a []string) {
	a = strings.Split(args, ",")
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
	proxyID, err = extractProxyID(resp)
	if err != nil {
		return "", err
	}
	return
}

func extractProxyID(resp zabbix.Response) (proxyID string, err error) {
	r, err := assertFirstResult(resp)
	if err != nil {
		return "", err
	}
	proxyID = r["proxyid"].(string)
	err = nil
	return
}

func assertFirstResult(resp zabbix.Response) (r map[string]interface{}, err error) {
	rr := resp.Result.([]interface{})
	r, ok := rr[0].(map[string]interface{})
	if !ok {
		err := errors.New("assertion error")
		return r, err
	}
	err = nil
	return
}

func setLoggerColog(debug bool) {
	colog.SetDefaultLevel(colog.LInfo)
	switch debug {
	case true:
		colog.SetMinLevel(colog.LDebug)
		colog.SetFormatter(&colog.StdFormatter{
			Colors: true,
			Flag:   log.Ldate | log.Ltime | log.Lshortfile,
		})
	default:
		colog.SetMinLevel(colog.LInfo)
		colog.SetFormatter(&colog.StdFormatter{
			Colors: true,
			Flag:   log.Ldate | log.Ltime,
		})
	}
	colog.Register()
	/*
			log.Printf("trace: this is a trace log.")
			log.Printf("debug: this is a debug log.")
			log.Printf("info: this is an info log.")
			log.Printf("warn: this is a warning log.")
			log.Printf("error: this is an error log.")
		    log.Printf("alert: this is an alert log.")
			log.Printf("this is a default level log.")
	*/
}

func outputTable(list []string, header string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{header})
	for _, v := range list {
		table.Append([]string{v})
	}
	table.Render()
}

func outputRaw(list []string) {
	for _, v := range list {
		fmt.Println(v)
	}
}
