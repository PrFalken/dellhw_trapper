package main

import (
	"strconv"
	"strings"
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

func dummy_report() (MultiDataPoint, error) {
	var md MultiDataPoint
	Add(&md, "hw.dummy", "0", TagSet{"test": "dummy"}, "Dummy description")
	return md, nil
}

func c_omreport_chassis() (MultiDataPoint, error) {
	var md MultiDataPoint
	readOmreport(func(fields []string) {
		if len(fields) != 2 || fields[0] == "SEVERITY" {
			return
		}
		component := strings.Replace(fields[1], " ", "_", -1)
		Add(&md, "hw.chassis", severity(fields[0]), TagSet{"component": component}, descDellHWChassis)
	}, "chassis")
	return md, nil
}

func c_omreport_system() (MultiDataPoint, error) {
	var md MultiDataPoint
	readOmreport(func(fields []string) {
		if len(fields) != 2 || fields[0] == "SEVERITY" {
			return
		}
		component := strings.Replace(fields[1], " ", "_", -1)
		Add(&md, "hw.system", severity(fields[0]), TagSet{"component": component}, descDellHWSystem)
	}, "system")
	return md, nil
}

func c_omreport_storage_enclosure() (MultiDataPoint, error) {
	var md MultiDataPoint
	readOmreport(func(fields []string) {
		if len(fields) < 3 || fields[0] == "ID" {
			return
		}
		id := strings.Replace(fields[0], ":", "_", -1)
		Add(&md, "hw.storage.enclosure", fields[1], TagSet{"id": id}, descDellHWStorageEnc)
	}, "storage", "enclosure")
	return md, nil
}

func c_omreport_storage_vdisk() (MultiDataPoint, error) {
	var md MultiDataPoint
	readOmreport(func(fields []string) {
		if len(fields) < 3 || fields[0] == "ID" {
			return
		}
		id := strings.Replace(fields[0], ":", "_", -1)
		Add(&md, "hw.storage.vdisk", fields[1], TagSet{"id": id}, descDellHWVDisk)
	}, "storage", "vdisk")
	return md, nil
}

func c_omreport_ps() (MultiDataPoint, error) {
	var md MultiDataPoint
	readOmreport(func(fields []string) {
		if len(fields) < 3 || fields[0] == "Index" {
			return
		}
		id := strings.Replace(fields[0], ":", "_", -1)
		ts := TagSet{"id": id}
		Add(&md, "hw.ps", fields[1], ts, descDellHWPS)
		if len(fields) < 6 {
			return
		}
		if fields[4] != "" {
			Add(&md, "hw.rated_input_wattage", fields[4], nil, descDellHWPS)
		}
		if fields[5] != "" {
			Add(&md, "hw.rated_output_wattage", fields[5], nil, descDellHWPS)
		}
	}, "chassis", "pwrsupplies")
	return md, nil
}

func c_omreport_ps_amps_sysboard_pwr() (MultiDataPoint, error) {
	var md MultiDataPoint
	readOmreport(func(fields []string) {
		if len(fields) == 2 && strings.Contains(fields[0], "Current") {
			i_fields := strings.Split(fields[0], "Current")
			v_fields := strings.Fields(fields[1])
			if len(i_fields) < 2 && len(v_fields) < 2 {
				return
			}
			id := strings.Replace(i_fields[0], " ", "", -1)
			Add(&md, "hw.chassis.current.reading", v_fields[0], TagSet{"id": id}, descDellHWCurrent)
		} else if len(fields) == 6 && (fields[2] == "System Board Pwr Consumption" || fields[2] == "System Board System Level") {
			v_fields := strings.Fields(fields[3])
			warn_fields := strings.Fields(fields[4])
			fail_fields := strings.Fields(fields[5])
			if len(v_fields) < 2 || len(warn_fields) < 2 || len(fail_fields) < 2 {
				return
			}
			Add(&md, "hw.chassis.power.reading", v_fields[0], nil, descDellHWPower)
			Add(&md, "hw.chassis.power.warn_level", warn_fields[0], nil, descDellHWPowerThreshold)
			Add(&md, "hw.chassis.power.fail_level", fail_fields[0], nil, descDellHWPowerThreshold)
		}
	}, "chassis", "pwrmonitoring")
	return md, nil
}

func c_omreport_storage_battery() (MultiDataPoint, error) {
	var md MultiDataPoint
	readOmreport(func(fields []string) {
		if len(fields) < 3 || fields[0] == "ID" {
			return
		}
		id := strings.Replace(fields[0], ":", "_", -1)
		Add(&md, "hw.storage.battery", fields[1], TagSet{"id": id}, descDellHWStorageBattery)
	}, "storage", "battery")
	return md, nil
}

func c_omreport_storage_controller() (MultiDataPoint, error) {
	var md MultiDataPoint
	readOmreport(func(fields []string) {
		if len(fields) < 3 || fields[0] == "ID" {
			return
		}
		c_omreport_storage_pdisk(fields[0], &md)
		id := strings.Replace(fields[0], ":", "_", -1)
		ts := TagSet{"id": id}
		Add(&md, "hw.storage.controller", fields[1], ts, descDellHWStorageCtl)
	}, "storage", "controller")
	return md, nil
}

// c_omreport_storage_pdisk is called from the controller func, since it needs the encapsulating id.
func c_omreport_storage_pdisk(id string, md *MultiDataPoint) {
	readOmreport(func(fields []string) {
		if len(fields) < 3 || fields[0] == "ID" {
			return
		}
		//Need to find out what the various ID formats might be
		id := strings.Replace(fields[0], ":", "_", -1)
		ts := TagSet{"id": id}
		Add(md, "hw.storage.pdisk", fields[1], ts, descDellHWPDisk)
	}, "storage", "pdisk", "controller="+id)
}

func c_omreport_processors() (MultiDataPoint, error) {
	var md MultiDataPoint
	readOmreport(func(fields []string) {
		if len(fields) != 8 {
			return
		}
		if _, err := strconv.Atoi(fields[0]); err != nil {
			return
		}
		ts := TagSet{"name": replace(fields[2])}
		Add(&md, "hw.chassis.processor", fields[1], ts, descDellHWCPU)
		AddMeta("", ts, "processor", clean(fields[3], fields[4]), true)
	}, "chassis", "processors")
	return md, nil
}

func c_omreport_fans() (MultiDataPoint, error) {
	var md MultiDataPoint
	readOmreport(func(fields []string) {
		if len(fields) != 8 {
			return
		}
		if _, err := strconv.Atoi(fields[0]); err != nil {
			return
		}
		ts := TagSet{"name": replace(fields[2])}
		Add(&md, "hw.chassis.fan", fields[1], ts, descDellHWFan)
		fs := strings.Fields(fields[3])
		if len(fs) == 2 && fs[1] == "RPM" {
			i, err := strconv.Atoi(fs[0])
			if err == nil {
				Add(&md, "hw.chassis.fan.reading", i, ts, descDellHWFanSpeed)
			}
		}
	}, "chassis", "fans")
	return md, nil
}

func c_omreport_memory() (MultiDataPoint, error) {
	var md MultiDataPoint
	readOmreport(func(fields []string) {
		if len(fields) != 5 {
			return
		}
		if _, err := strconv.Atoi(fields[0]); err != nil {
			return
		}
		ts := TagSet{"name": replace(fields[2])}
		Add(&md, "hw.chassis.memory", fields[1], ts, descDellHWMemory)
		AddMeta("", ts, "memory", clean(fields[4]), true)
	}, "chassis", "memory")
	return md, nil
}

func c_omreport_temps() (MultiDataPoint, error) {
	var md MultiDataPoint
	readOmreport(func(fields []string) {
		if len(fields) != 8 {
			return
		}
		if _, err := strconv.Atoi(fields[0]); err != nil {
			return
		}
		ts := TagSet{"name": replace(fields[2])}
		Add(&md, "hw.chassis.temps", fields[1], ts, descDellHWTemp)
		fs := strings.Fields(fields[3])
		if len(fs) == 2 && fs[1] == "C" {
			i, err := strconv.ParseFloat(fs[0], 64)
			if err == nil {
				Add(&md, "hw.chassis.temps.reading", i, ts, descDellHWTempReadings)
			}
		}
	}, "chassis", "temps")
	return md, nil
}

func c_omreport_volts() (MultiDataPoint, error) {
	var md MultiDataPoint
	readOmreport(func(fields []string) {
		if len(fields) != 8 {
			return
		}
		if _, err := strconv.Atoi(fields[0]); err != nil {
			return
		}
		ts := TagSet{"name": replace(fields[2])}
		Add(&md, "hw.chassis.volts", fields[1], ts, descDellHWVolt)
		if i, err := extract(fields[3], "V"); err == nil {
			Add(&md, "hw.chassis.volts.reading", i, ts, descDellHWVoltReadings)
		}
	}, "chassis", "volts")
	return md, nil
}
