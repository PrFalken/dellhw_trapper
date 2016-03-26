# dellhw_trapper

## Zabbix exporter for Dell Hardware components

*Supports Dell OMSA 7.4*

	Usage:
	  dellhw_trapper [flags]
	  dellhw_trapper [command]
	
	Available Commands:
	  version     Print the version number of hardware_exporter
	  help        Help about any command
	
	Flags:
	  -c, --collect="chassis,fans,memory,processors,ps,ps_amps_sysboard_pwr,storage_battery,storage_enclosure,storage_controller,storage_vdisk,system,temps,volts": Comma-separated list of collectors to use.
	      --discovery[=false]: Perform Zabbix low level discovery on hardware elements
	  -h, --help[=false]: help for dellhw_trapper
	  -L, --loglevel="info": Set log level
	  -n, --namespace="": Discovery key
	      --update-items[=false]: Get & send items to Zabbix. This is the default behaviour
	  -f, --zabbix-from="lucky.local": Send to Zabbix from this host name. You can also set HOSTNAME and DOMAINNAME environment variables.
	  -p, --zabbix-port="10051": Zabbix server port
	  -z, --zabbix-server="localhost": Zabbix server hostname or address
	
	
	Use "dellhw_trapper [command] --help" for more information about a command.
