package command

import (
	"errors"
	"log"

	"github.com/AlekSi/zabbix"
	"github.com/codegangsta/cli"
)

func CmdCreate(c *cli.Context) (err error) {
	// set logger
	setLoggerColog(c.GlobalBool("debug"))

	// Required args check
	m := map[string]string{
		"hostname":   c.String("hostname"),
		"host-group": c.String("host-group"),
	}
	checkRequiredStringArgs(m)

	z := newZabbixctl(c)
	if err = z.login(); err != nil {
		log.Printf("error: %v", err)
		return
	}

	exist, err := z.hostExists(c.String("hostname"))
	if err != nil {
		log.Printf("error: %v", err)
		return
	} else if exist == true {
		err = errors.New("Host is already exist.")
		log.Printf("error: %v", err)
		return err
	}

	_, err = z.hostCreate(
		c.String("hostname"),
		c.String("host-group"),
		c.String("templates"),
		c.String("ipaddress"),
		c.String("dnsname"),
		c.String("port"),
		c.String("proxy"),
		c.Bool("use-ip"),
	)
	if err != nil {
		log.Printf("error: %v", err)
	} else {
		log.Printf("info: Host is created.")
	}
	return
}

func (z *zabbixcli) hostCreate(hostname, hostGroup, templates, ipaddress, dnsname, port, proxy string, useIp bool) (resp zabbix.Response, err error) {
	zabbixParams := z.buildZabbixParams(hostname, hostGroup, templates, ipaddress, dnsname, port, proxy, useIp)

	resp, err = z.Api.Call("host.create", zabbixParams)
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

func (z *zabbixcli) buildZabbixParams(hostname, hostGroup, templates, ipaddress, dnsname, port, proxy string, useIp bool) (zabbixParams zabbix.Params) {
	useInterface := whichInterfaceToBeUsed(useIp, ipaddress, dnsname)
	var hostGroupIDs []map[string]string
	hostGroupList := splitArgs(hostGroup)
	for _, hostGroup := range hostGroupList {
		hostGroupID, err := z.hostGroupIdGet(hostGroup)
		if err != nil {
			log.Fatal(err)
		}
		hostGroupIDs = append(hostGroupIDs, map[string]string{"groupid": hostGroupID})
	}

	// テンプレート使ってるかチェック
	useTemplate := true
	var Templates []map[string]string
	Templates = func(templates string) []map[string]string {
		if templates == "" {
			useTemplate = false
			return Templates
		} else {
			templateList := splitArgs(templates)
			for _, t := range templateList {
				templateID, err := z.templateIdGet(t)
				if err != nil {
					log.Fatal(err)
				}
				Templates = append(Templates, map[string]string{"templateid": templateID})
			}
			return Templates
		}
	}(templates)

	// proxy使ってるかチェック
	useProxy := true
	proxyID := func(proxy string) string {
		if ifEmptyString(proxy) {
			useProxy = false
			return ""
		} else {
			proxyID, err := z.proxyIdGet(proxy)
			if err != nil {
				log.Fatal(err)
			}
			return proxyID
		}
	}(proxy)

	return func(hostname, ipaddress, dnsname, port, proxyID string,
		useInterface int,
		Templates, hostGroupIDs []map[string]string,
		useTemplate, useProxy bool) zabbix.Params {
		// build zabbix.Params
		if useProxy {
			if useTemplate {
				// use proxy, use template
				return zabbix.Params{
					"host": hostname,
					"interfaces": []interface{}{
						map[string]interface{}{
							"type":  1,
							"main":  1,
							"useip": useInterface,
							"ip":    ipaddress,
							"dns":   dnsname,
							"port":  port,
						},
					},
					"groups":       hostGroupIDs,
					"templates":    Templates,
					"proxy_hostid": proxyID,
				}
			} else {
				// use only proxy
				return zabbix.Params{
					"host": hostname,
					"interfaces": []interface{}{
						map[string]interface{}{
							"type":  1,
							"main":  1,
							"useip": useInterface,
							"ip":    ipaddress,
							"dns":   dnsname,
							"port":  port,
						},
					},
					"groups":       hostGroupIDs,
					"proxy_hostid": proxyID,
				}
			}
		} else {
			if useTemplate {
				// use only template
				return zabbix.Params{
					"host": hostname,
					"interfaces": []interface{}{
						map[string]interface{}{
							"type":  1,
							"main":  1,
							"useip": useInterface,
							"ip":    ipaddress,
							"dns":   dnsname,
							"port":  port,
						},
					},
					"groups":    hostGroupIDs,
					"templates": Templates,
				}
			} else {
				// not use proxy and template
				return zabbix.Params{
					"host": hostname,
					"interfaces": []interface{}{
						map[string]interface{}{
							"type":  1,
							"main":  1,
							"useip": useInterface,
							"ip":    ipaddress,
							"dns":   dnsname,
							"port":  port,
						},
					},
					"groups": hostGroupIDs,
				}
			}
		}
	}(hostname, ipaddress, dnsname, port, proxyID, useInterface, Templates, hostGroupIDs, useTemplate, useProxy)
}
