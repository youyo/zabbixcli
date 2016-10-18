package command

import (
	"log"
	"os"

	"github.com/AlekSi/zabbix"
	"github.com/codegangsta/cli"
	"github.com/olekukonko/tablewriter"
)

func CmdTemplateList(c *cli.Context) (err error) {
	// set logger
	setLoggerColog(c.GlobalBool("debug"))

	z := newZabbixctl(c)
	if err = z.login(); err != nil {
		log.Printf("error: %v", err)
		return
	}
	templates, err := z.templateGet()
	if err != nil {
		log.Printf("error: %v", err)
		return
	}

	// select output format
	switch c.Bool("raw") {
	case true:
		outputRaw(templates)
	default:
		outputTable(templates, "Templates")
	}
	return
}

func (z *zabbixctl) templateGet() (templates []string, err error) {
	resp, err := z.Api.Call("template.get", zabbix.Params{
		"output":    "extend",
		"sortfield": "name",
	})
	if err != nil {
		log.Printf("error: %v", err)
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
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Templates"})
	for _, v := range templates {
		table.Append([]string{v})
	}
	table.Render()
}
