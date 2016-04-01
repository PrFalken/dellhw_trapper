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


Example of discovered metrics on a Dell PowerEdge R630

	dell.hardware.chassis.current[reading]:0.2
	dell.hardware.chassis.power.fail[level]:644
	dell.hardware.chassis.power.warn[level]:588
	dell.hardware.chassis.power[reading]:126
	dell.hardware.chassis.temps[number]:4
	dell.hardware.chassis.temps[reading]:37.0
	dell.hardware.chassis.temps[status]:0
	dell.hardware.chassis.temps[status_sum]:0
	dell.hardware.chassis.volts[number]:32
	dell.hardware.chassis.volts[reading]:230
	dell.hardware.chassis.volts[status]:0
	dell.hardware.chassis.volts[status_sum]:0
	dell.hardware.chassis[number]:10
	dell.hardware.chassis[status]:1
	dell.hardware.chassis[status_sum]:1
	dell.hardware.fan[System Board Fan1A,speed]:4920
	dell.hardware.fan[System Board Fan1A,status]:0
	dell.hardware.fan[System Board Fan1B,speed]:4560
	dell.hardware.fan[System Board Fan1B,status]:0
	dell.hardware.fan[System Board Fan2A,speed]:4800
	dell.hardware.fan[System Board Fan2A,status]:0
	dell.hardware.fan[System Board Fan2B,speed]:4560
	dell.hardware.fan[System Board Fan2B,status]:0
	dell.hardware.fan[System Board Fan3A,speed]:4920
	dell.hardware.fan[System Board Fan3A,status]:0
	dell.hardware.fan[System Board Fan3B,speed]:4560
	dell.hardware.fan[System Board Fan3B,status]:0
	dell.hardware.fan[System Board Fan4A,speed]:4920
	dell.hardware.fan[System Board Fan4A,status]:0
	dell.hardware.fan[System Board Fan4B,speed]:4560
	dell.hardware.fan[System Board Fan4B,status]:0
	dell.hardware.fan[System Board Fan5A,speed]:4920
	dell.hardware.fan[System Board Fan5A,status]:0
	dell.hardware.fan[System Board Fan5B,speed]:4560
	dell.hardware.fan[System Board Fan5B,status]:0
	dell.hardware.fan[System Board Fan6A,speed]:4920
	dell.hardware.fan[System Board Fan6A,status]:0
	dell.hardware.fan[System Board Fan6B,speed]:4560
	dell.hardware.fan[System Board Fan6B,status]:0
	dell.hardware.fan[System Board Fan7A,speed]:5160
	dell.hardware.fan[System Board Fan7A,status]:0
	dell.hardware.fan[System Board Fan7B,speed]:4800
	dell.hardware.fan[System Board Fan7B,status]:0
	dell.hardware.fan[number]:14
	dell.hardware.fan[status_sum]:0
	dell.hardware.memory[A1,status]:0
	dell.hardware.memory[A2,status]:0
	dell.hardware.memory[A3,status]:0
	dell.hardware.memory[A4,status]:0
	dell.hardware.memory[B1,status]:0
	dell.hardware.memory[B2,status]:0
	dell.hardware.memory[B3,status]:0
	dell.hardware.memory[B4,status]:0
	dell.hardware.memory[number]:8
	dell.hardware.memory[status_sum]:0
	dell.hardware.power[0,input_watts]:594
	dell.hardware.power[0,output_watts]:495
	dell.hardware.power[0,status]:0
	dell.hardware.power[1,input_watts]:594
	dell.hardware.power[1,output_watts]:495
	dell.hardware.power[1,status]:0
	dell.hardware.power[number]:2
	dell.hardware.power[status_sum]:0
	dell.hardware.processors[CPU1,status]:0
	dell.hardware.processors[CPU2,status]:0
	dell.hardware.processors[number]:2
	dell.hardware.processors[status_sum]:0
	dell.hardware.raid.controller[0,status]:0
	dell.hardware.raid.controller[number]:1
	dell.hardware.raid.controller[status_sum]:0
	dell.hardware.raid.logicaldrive[0,status]:0
	dell.hardware.raid.logicaldrive[1,status]:0
	dell.hardware.raid.logicaldrive[number]:2
	dell.hardware.raid.logicaldrive[status_sum]:0
	dell.hardware.raid.physicaldrive[0_1_0,status]:0
	dell.hardware.raid.physicaldrive[0_1_1,status]:0
	dell.hardware.raid.physicaldrive[0_1_2,status]:0
	dell.hardware.raid.physicaldrive[0_1_3,status]:0
	dell.hardware.raid.physicaldrive[0_1_4,status]:0
	dell.hardware.raid.physicaldrive[0_1_5,status]:0
	dell.hardware.raid.physicaldrive[number]:6
	dell.hardware.raid.physicaldrive[status_sum]:0
	dell.hardware.storage.battery[number]:1
	dell.hardware.storage.battery[status]:0
	dell.hardware.storage.battery[status_sum]:0
	dell.hardware.storage.enclosure[number]:1
	dell.hardware.storage.enclosure[status]:0
	dell.hardware.storage.enclosure[status_sum]:0
	dell.hardware.system[number]:1
	dell.hardware.system[status]:0
	dell.hardware.system[status_sum]:0