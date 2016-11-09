# zabbixcli

[![wercker status](https://app.wercker.com/status/2357506da6c5b2d9c86652321466787e/s/master "wercker status")](https://app.wercker.com/project/byKey/2357506da6c5b2d9c86652321466787e)

## Description

For control the host of zabbix.

## Usage

```
$ zabbixcli -h
NAME:
   zabbixcli - For control the host of zabbix.

USAGE:
   zabbixcli [global options] command [command options] [arguments...]

VERSION:
   0.1.1

AUTHOR(S):
   youyo

COMMANDS:
     help, h  Shows a list of commands or help for one command

   Host management:
     create   Create host
     enable   Enable host
     disable  Disable host
     delete   Delete host

   Lists:
     proxy-list      Shows a list of zabbix-proxies
     template-list   Shows a list of templates
     hostgroup-list  Shows a list of host-groups
     search          Search hosts

GLOBAL OPTIONS:
   --debug           Set LogLevel Debug.
   --username value  set username. [$ZABBIXCTL_USERNAME]
   --password value  set password. [$ZABBIXCTL_PASSWORD]
   --url value       set url. [$ZABBIXCTL_URL]
   --help, -h        show help
   --version, -v     print the version
```

## Install

To install,

```bash
$ wget https://github.com/youyo/zabbixcli/releases/download/$latest_version/zabbixcli_linux_amd64.zip
$ unzip zabbixcli_linux_amd64.zip -d /usr/local/bin/
```

## Contribution

1. Fork ([https://github.com/youyo/zabbixcli/fork](https://github.com/youyo/zabbixcli/fork))
1. Run `make setup && make deps`
1. Create a feature branch
1. Commit your changes
1. Rebase your local changes against the master branch
1. Run test suite with the `make test` command and confirm that it passes
1. Run `gofmt -s`
1. Create a new Pull Request

## Author

[youyo](https://github.com/youyo)
