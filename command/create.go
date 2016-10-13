package command

import (
	"errors"
	"log"

	"github.com/AlekSi/zabbix"
	"github.com/codegangsta/cli"
)

func CmdCreate(c *cli.Context) (err error) {
	// Required args check
	m := map[string]string{
		"hostname":   c.String("hostname"),
		"host-group": c.String("host-group"),
	}
	checkRequiredStringArgs(m)

	hostname := c.String("hostname")
	hostGroup := c.String("host-group")
	templates := c.String("templates")
	ipaddress := c.String("ipaddress")
	dnsname := c.String("dnsname")
	proxy := c.String("proxy")
	useIp := c.Bool("use-ip")

	z := newZabbixctl(c)
	if err = z.login(); err != nil {
		return
	}

	exist, err := z.hostExists(hostname)
	if err != nil {
		return
	} else if exist == true {
		return errors.New("Host is already exist.")
	}

	_, err = z.hostCreate(hostname, hostGroup, templates, ipaddress, dnsname, proxy, useIp)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func (z *zabbixctl) hostCreate(hostname, hostGroup, templates, ipaddress, dnsname, proxy string, useIp bool) (resp zabbix.Response, err error) {
	useInterface := whichInterfaceToBeUsed(useIp, ipaddress, dnsname)
	groupID, err := z.hostGroupIdGet(hostGroup)
	if err != nil {
		log.Fatal(err)
	}
	var Templates []map[string]string
	templateList := splitTemplates(templates)
	for _, t := range templateList {
		templateID, err := z.templateIdGet(t)
		if err != nil {
			log.Fatal(err)
		}
		Templates = append(Templates, map[string]string{"templateid": templateID})
	}
	proxyID, err := z.proxyIdGet(proxy)
	if err != nil {
		log.Fatal(err)
	}

	resp, err = z.Api.Call("host.create", zabbix.Params{
		"host": hostname,
		"interfaces": []interface{}{
			map[string]interface{}{
				"type":  1,
				"main":  1,
				"useip": useInterface,
				"ip":    ipaddress,
				"dns":   dnsname,
				"port":  "10050",
			},
		},
		"groups": []interface{}{
			map[string]string{
				"groupid": groupID,
			},
		},
		"templates":    Templates,
		"proxy_hostid": proxyID,
	})
	return
}

func whichInterfaceToBeUsed(useIp bool, ipaddress, dnsname string) int {
	switch useIp {
	case true:
		checkRequiredStringArgs(map[string]string{"ipaddress": ipaddress})
		return UseIp
	default:
		checkRequiredStringArgs(map[string]string{"dnsname": dnsname})
		return UseDns
	}
}
