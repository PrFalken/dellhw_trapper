package main

import (
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	descDellHWChassis        = "Overall status of chassis components."
	descDellHWSystem         = "Overall status of system components."
	descDellHWStorageEnc     = "Overall status of storage enclosures."
	descDellHWVDisk          = "Overall status of virtual disks."
	descDellHWPS             = "Overall status of power supplies."
	descDellHWCurrent        = "Amps used per power supply."
	descDellHWPower          = "System board power usage."
	descDellHWPowerThreshold = "The warning and failure levels set on the device for system board power usage."
	descDellHWStorageBattery = "Status of storage controller backup batteries."
	descDellHWStorageCtl     = "Overall status of storage controllers."
	descDellHWPDisk          = "Overall status of physical disks."
	descDellHWCPU            = "Overall status of CPUs."
	descDellHWFan            = "Overall status of system fans."
	descDellHWFanSpeed       = "System fan speed."
	descDellHWMemory         = "System RAM DIMM status."
	descDellHWTemp           = "Overall status of system temperature readings."
	descDellHWTempReadings   = "System temperature readings."
	descDellHWVolt           = "Overall status of power supply volt readings."
	descDellHWVoltReadings   = "Volts used per power supply."
)

func readOmreport(f func([]string), args ...string) {
	args = append(args, "-fmt", "ssv")
	_ = readCommand(func(line string) error {
		sp := strings.Split(line, ";")
		for i, s := range sp {
			sp[i] = clean(s)
		}
		f(sp)
		return nil
	}, "/opt/dell/srvadmin/bin/omreport", args...)
}

func dummy_report() error {
	Add("dummy", "1", prometheus.Labels{"test": "dummy"}, "Dummy description")
	return nil
}

func c_omreport_chassis() error {
	readOmreport(func(fields []string) {
		if len(fields) != 2 || fields[0] == "SEVERITY" {
			return
		}
		component := strings.Replace(fields[1], " ", "_", -1)
		Add("chassis", severity(fields[0]), prometheus.Labels{"component": component}, descDellHWChassis)
	}, "chassis")
	return nil
}

func c_omreport_system() error {
	readOmreport(func(fields []string) {
		if len(fields) != 2 || fields[0] == "SEVERITY" {
			return
		}
		component := strings.Replace(fields[1], " ", "_", -1)
		Add("system", severity(fields[0]), prometheus.Labels{"component": component}, descDellHWSystem)
	}, "system")
	return nil
}

func c_omreport_storage_enclosure() error {
	readOmreport(func(fields []string) {
		if len(fields) < 3 || fields[0] == "ID" {
			return
		}
		id := strings.Replace(fields[0], ":", "_", -1)
		Add("storage_enclosure", severity(fields[1]), prometheus.Labels{"id": id}, descDellHWStorageEnc)
	}, "storage", "enclosure")
	return nil
}

func c_omreport_storage_vdisk() error {
	readOmreport(func(fields []string) {
		if len(fields) < 3 || fields[0] == "ID" {
			return
		}
		id := strings.Replace(fields[0], ":", "_", -1)
		Add("storage_vdisk", severity(fields[1]), prometheus.Labels{"id": id}, descDellHWVDisk)
	}, "storage", "vdisk")
	return nil
}

func c_omreport_ps() error {
	readOmreport(func(fields []string) {
		if len(fields) < 3 || fields[0] == "Index" {
			return
		}
		id := strings.Replace(fields[0], ":", "_", -1)
		ts := prometheus.Labels{"id": id}
		Add("ps", severity(fields[1]), ts, descDellHWPS)
		if len(fields) < 6 {
			return
		}
		if fields[4] != "" {
			iWattage, err := extract(fields[4], "W")
			if err == nil {
				Add("rated_input_wattage", iWattage, ts, descDellHWPS)
			}
		}
		if fields[5] != "" {
			oWattage, err := extract(fields[5], "W")
			if err == nil {
				Add("rated_output_wattage", oWattage, ts, descDellHWPS)
			}
		}
	}, "chassis", "pwrsupplies")
	return nil
}

func c_omreport_ps_amps_sysboard_pwr() error {
	readOmreport(func(fields []string) {
		if len(fields) == 2 && strings.Contains(fields[0], "Current") {
			iFields := strings.Split(fields[0], "Current")
			vFields := strings.Fields(fields[1])
			if len(iFields) < 2 && len(vFields) < 2 {
				return
			}
			id := strings.Replace(iFields[0], " ", "", -1)
			Add("chassis_current_reading", vFields[0], prometheus.Labels{"id": id}, descDellHWCurrent)
		} else if len(fields) == 6 && (fields[2] == "System Board Pwr Consumption" || fields[2] == "System Board System Level") {
			vFields := strings.Fields(fields[3])
			warnFields := strings.Fields(fields[4])
			failFields := strings.Fields(fields[5])
			if len(vFields) < 2 || len(warnFields) < 2 || len(failFields) < 2 {
				return
			}
			Add("chassis_power_reading", vFields[0], nil, descDellHWPower)
			Add("chassis_power_warn_level", warnFields[0], nil, descDellHWPowerThreshold)
			Add("chassis_power_fail_level", failFields[0], nil, descDellHWPowerThreshold)
		}
	}, "chassis", "pwrmonitoring")
	return nil
}

func c_omreport_storage_battery() error {
	readOmreport(func(fields []string) {
		if len(fields) < 3 || fields[0] == "ID" {
			return
		}
		id := strings.Replace(fields[0], ":", "_", -1)
		Add("storage_battery", severity(fields[1]), prometheus.Labels{"id": id}, descDellHWStorageBattery)
	}, "storage", "battery")
	return nil
}

func c_omreport_storage_controller() error {
	readOmreport(func(fields []string) {
		if len(fields) < 3 || fields[0] == "ID" {
			return
		}
		c_omreport_storage_pdisk(fields[0])
		id := strings.Replace(fields[0], ":", "_", -1)
		ts := prometheus.Labels{"id": id}
		Add("storage_controller", fields[1], ts, descDellHWStorageCtl)
	}, "storage", "controller")
	return nil
}

// c_omreport_storage_pdisk is called from the controller func, since it needs the encapsulating id.
func c_omreport_storage_pdisk(id string) {
	readOmreport(func(fields []string) {
		if len(fields) < 3 || fields[0] == "ID" {
			return
		}
		//Need to find out what the various ID formats might be
		id := strings.Replace(fields[0], ":", "_", -1)
		ts := prometheus.Labels{"id": id}
		Add("storage_pdisk", fields[1], ts, descDellHWPDisk)
	}, "storage", "pdisk", "controller="+id)
}

func c_omreport_processors() error {
	readOmreport(func(fields []string) {
		if len(fields) != 8 {
			return
		}
		if _, err := strconv.Atoi(fields[0]); err != nil {
			return
		}
		ts := prometheus.Labels{"name": replace(fields[2])}
		Add("chassis_processor", severity(fields[1]), ts, descDellHWCPU)
	}, "chassis", "processors")
	return nil
}

func c_omreport_fans() error {
	readOmreport(func(fields []string) {
		if len(fields) != 8 {
			return
		}
		if _, err := strconv.Atoi(fields[0]); err != nil {
			return
		}
		ts := prometheus.Labels{"name": replace(fields[2])}
		Add("chassis_fan", fields[1], ts, descDellHWFan)
		fs := strings.Fields(fields[3])
		if len(fs) == 2 && fs[1] == "RPM" {
			Add("chassis_fan_reading", fs[0], ts, descDellHWFanSpeed)
		}
	}, "chassis", "fans")
	return nil
}

func c_omreport_memory() error {
	readOmreport(func(fields []string) {
		if len(fields) != 5 {
			return
		}
		if _, err := strconv.Atoi(fields[0]); err != nil {
			return
		}
		ts := prometheus.Labels{"name": replace(fields[2])}
		Add("chassis_memory", severity(fields[1]), ts, descDellHWMemory)
	}, "chassis", "memory")
	return nil
}

func c_omreport_temps() error {
	readOmreport(func(fields []string) {
		if len(fields) != 8 {
			return
		}
		if _, err := strconv.Atoi(fields[0]); err != nil {
			return
		}
		ts := prometheus.Labels{"name": replace(fields[2])}
		Add("chassis_temps", severity(fields[1]), ts, descDellHWTemp)
		fs := strings.Fields(fields[3])
		if len(fs) == 2 && fs[1] == "C" {
			Add("chassis_temps_reading", fs[0], ts, descDellHWTempReadings)
		}
	}, "chassis", "temps")
	return nil
}

func c_omreport_volts() error {
	readOmreport(func(fields []string) {
		if len(fields) != 8 {
			return
		}
		if _, err := strconv.Atoi(fields[0]); err != nil {
			return
		}
		ts := prometheus.Labels{"name": replace(fields[2])}
		Add("chassis_volts", severity(fields[1]), ts, descDellHWVolt)
		if i, err := extract(fields[3], "V"); err == nil {
			Add("chassis_volts_reading", i, ts, descDellHWVoltReadings)
		}
	}, "chassis", "volts")
	return nil
}
