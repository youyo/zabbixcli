package command

import (
	"fmt"

	"github.com/AlekSi/zabbix"
	"github.com/codegangsta/cli"
)

func CmdTemplateList(c *cli.Context) (err error) {
	z := newZabbixctl(c)
	if err = z.login(); err != nil {
		return
	}
	templates, err := z.templateGet()
	outputTemplateList(templates)
	return
}

func (z *zabbixctl) templateGet() (templates []string, err error) {
	resp, err := z.Api.Call("template.get", zabbix.Params{
		"output":    "extend",
		"sortfield": "name",
	})
	if err != nil {
		return nil, err
	}
	templates, _ = extractTemplateName(resp)
	return
}

func extractTemplateName(resp zabbix.Response) (templates []string, err error) {
	rr := resp.Result.([]interface{})
	for _, r := range rr {
		r := r.(map[string]interface{})
		templates = append(templates, r["name"].(string))
	}
	return templates, nil
}

func outputTemplateList(templates []string) {
	for _, v := range templates {
		fmt.Println(v)
	}
}
