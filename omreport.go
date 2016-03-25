package main

import (
	"strconv"
	"strings"

	log "github.com/Sirupsen/logrus"
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

var (
	collectors = map[string]collector{
		"dummy":      collector{F: dummyReport},
		"chassis":    collector{F: omreportChassis},
		"fans":       collector{F: omreportFans},
		"memory":     collector{F: omreportMemory},
		"processors": collector{F: omreportProcessors},
		"ps":         collector{F: omreportPs},
		"ps_amps_sysboard_pwr": collector{F: omreportPsAmpsSysboardPwr},
		"storage_battery":      collector{F: omreportStorageBattery},
		"storage_controller":   collector{F: omreportStorageController},
		"storage_enclosure":    collector{F: omreportStorageEnclosure},
		"storage_vdisk":        collector{F: omreportStorageVdisk},
		"system":               collector{F: omreportSystem},
		"temps":                collector{F: omreportTemps},
		"volts":                collector{F: omreportVolts},
	}
)

type collector struct {
	F func(omReporter) error
}

type labels map[string]string

func collect(collectors map[string]collector) error {
	for _, name := range strings.Split(enabledCollectors, ",") {
		collector := collectors[name]
		log.Debug("Running collector ", name)
		om := newOmReport()
		err := collector.F(om)
		if err != nil {
			log.Error("Collector", name, "failed to run")
			return err
		}
	}
	return nil
}

type omReport struct{}

func newOmReport() *omReport {
	return &omReport{}
}

func (o *omReport) Report(f func([]string), args ...string) {
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

type omReporter interface {
	Report(f func([]string), args ...string)
}

func add(name string, value string, t labels, desc string) {

	cache.Lock.Lock()
	defer cache.Lock.Unlock()
	metric := newZabbixItem(name, t, value, desc)
	cache.metrics[name] = *metric

}

func dummyReport(om omReporter) error {
	add("dummy", "1", labels{"#{FUNKY}": "lolilol"}, "Dummy description")
	return nil
}

func omreportChassis(om omReporter) error {
	om.Report(func(fields []string) {
		if len(fields) != 2 || fields[0] == "SEVERITY" {
			return
		}
		component := strings.Replace(fields[1], " ", "_", -1)
		add("chassis", severity(fields[0]), labels{"component": component}, descDellHWChassis)
	}, "chassis")
	return nil
}

func omreportSystem(om omReporter) error {
	om.Report(func(fields []string) {
		if len(fields) != 2 || fields[0] == "SEVERITY" {
			return
		}
		component := strings.Replace(fields[1], " ", "_", -1)
		add("system", severity(fields[0]), labels{"component": component}, descDellHWSystem)
	}, "system")
	return nil
}

func omreportStorageEnclosure(om omReporter) error {
	om.Report(func(fields []string) {
		if len(fields) < 3 || fields[0] == "ID" {
			return
		}
		id := strings.Replace(fields[0], ":", "_", -1)
		add("storage_enclosure", severity(fields[1]), labels{"id": id}, descDellHWStorageEnc)
	}, "storage", "enclosure")
	return nil
}

func omreportStorageVdisk(om omReporter) error {
	om.Report(func(fields []string) {
		if len(fields) < 3 || fields[0] == "ID" {
			return
		}
		id := strings.Replace(fields[0], ":", "_", -1)
		ts := labels{"{#LOGICALDRIVESLOT}": id}
		add("dell.hardware.raid.logicaldrive["+id+",status]", severity(fields[1]), ts, descDellHWVDisk)
	}, "storage", "vdisk")
	return nil
}

func omreportPs(om omReporter) error {
	om.Report(func(fields []string) {
		if len(fields) < 3 || fields[0] == "Index" {
			return
		}
		id := strings.Replace(fields[0], ":", "_", -1)
		ts := labels{"{#POWERSLOT}": id}
		add("dell.hardware.power["+id+",status]", severity(fields[1]), ts, descDellHWPS)
		if len(fields) < 6 {
			return
		}
		if fields[4] != "" {
			iWattage, err := extract(fields[4], "W")
			if err == nil {
				add("dell.hardware.power["+id+",input_watts]", iWattage, ts, descDellHWPS)
			}
		}
		if fields[5] != "" {
			oWattage, err := extract(fields[5], "W")
			if err == nil {
				add("dell.hardware.power["+id+",ouput_watts]", oWattage, ts, descDellHWPS)
			}
		}
	}, "chassis", "pwrsupplies")
	return nil
}

func omreportPsAmpsSysboardPwr(om omReporter) error {
	om.Report(func(fields []string) {
		if len(fields) == 2 && strings.Contains(fields[0], "Current") {
			iFields := strings.Split(fields[0], "Current")
			vFields := strings.Fields(fields[1])
			if len(iFields) < 2 && len(vFields) < 2 {
				return
			}
			id := strings.Replace(iFields[0], " ", "", -1)
			add("chassis_current_reading", vFields[0], labels{"id": id}, descDellHWCurrent)
		} else if len(fields) == 6 && (fields[2] == "System Board Pwr Consumption" || fields[2] == "System Board System Level") {
			vFields := strings.Fields(fields[3])
			warnFields := strings.Fields(fields[4])
			failFields := strings.Fields(fields[5])
			if len(vFields) < 2 || len(warnFields) < 2 || len(failFields) < 2 {
				return
			}
			add("chassis_power_reading", vFields[0], nil, descDellHWPower)
			add("chassis_power_warn_level", warnFields[0], nil, descDellHWPowerThreshold)
			add("chassis_power_fail_level", failFields[0], nil, descDellHWPowerThreshold)
		}
	}, "chassis", "pwrmonitoring")
	return nil
}

func omreportStorageBattery(om omReporter) error {
	om.Report(func(fields []string) {
		if len(fields) < 3 || fields[0] == "ID" {
			return
		}
		id := strings.Replace(fields[0], ":", "_", -1)
		add("storage_battery", severity(fields[1]), labels{"id": id}, descDellHWStorageBattery)
	}, "storage", "battery")
	return nil
}

func omreportStorageController(om omReporter) error {
	om.Report(func(fields []string) {
		if len(fields) < 3 || fields[0] == "ID" {
			return
		}
		omreportStoragePdisk(om, fields[0])
		id := strings.Replace(fields[0], ":", "_", -1)
		ts := labels{"{#CONTROLLERSLOT}": id}
		add("dell.hardware.raid.controller["+id+",controller_status]", severity(fields[1]), ts, descDellHWStorageCtl)
	}, "storage", "controller")
	return nil
}

// omreportStoragePdisk is called from the controller func, since it needs the encapsulating id.
func omreportStoragePdisk(om omReporter, id string) {
	om.Report(func(fields []string) {
		if len(fields) < 3 || fields[0] == "ID" {
			return
		}
		//Need to find out what the various ID formats might be
		diskID := strings.Replace(fields[0], ":", "_", -1)
		ts := labels{"{#PHYSICALDRIVESLOT}": diskID, "{#CONTROLLERSLOT}": id}
		add("dell.hardware.raid.physicaldrive["+diskID+",status]", severity(fields[1]), ts, descDellHWPDisk)
	}, "storage", "pdisk", "controller="+id)
}

func omreportProcessors(om omReporter) error {
	om.Report(func(fields []string) {
		if len(fields) != 8 {
			return
		}
		if _, err := strconv.Atoi(fields[0]); err != nil {
			return
		}
		pname := replace(fields[2])
		ts := labels{"{#PROCESSORNAME}": pname}
		add("dell.hardware.processors["+pname+",status]", severity(fields[1]), ts, descDellHWCPU)
	}, "chassis", "processors")
	return nil
}

func omreportFans(om omReporter) error {
	om.Report(func(fields []string) {
		if len(fields) != 8 {
			return
		}
		if _, err := strconv.Atoi(fields[0]); err != nil {
			return
		}
		ts := labels{"{#FANNAME}": replace(fields[2])}
		add("dell.hardware.fan["+replace(fields[2])+",status]", severity(fields[1]), ts, descDellHWFan)
		fs := strings.Fields(fields[3])
		if len(fs) == 2 && fs[1] == "RPM" {
			add("dell.hardware.fan["+replace(fields[2])+",speed]", fs[0], ts, descDellHWFanSpeed)
		}
	}, "chassis", "fans")
	return nil
}

func omreportMemory(om omReporter) error {
	om.Report(func(fields []string) {
		if len(fields) != 5 {
			return
		}
		if _, err := strconv.Atoi(fields[0]); err != nil {
			return
		}
		slot := replace(fields[2])
		ts := labels{"{#MEMORYSLOT}": slot}
		add("dell.hardware.memory["+slot+",status]", severity(fields[1]), ts, descDellHWMemory)
	}, "chassis", "memory")
	return nil
}

func omreportTemps(om omReporter) error {
	om.Report(func(fields []string) {
		if len(fields) != 8 {
			return
		}
		if _, err := strconv.Atoi(fields[0]); err != nil {
			return
		}
		ts := labels{"name": replace(fields[2])}
		add("chassis_temps", severity(fields[1]), ts, descDellHWTemp)
		fs := strings.Fields(fields[3])
		if len(fs) == 2 && fs[1] == "C" {
			add("chassis_temps_reading", fs[0], ts, descDellHWTempReadings)
		}
	}, "chassis", "temps")
	return nil
}

func omreportVolts(om omReporter) error {
	om.Report(func(fields []string) {
		if len(fields) != 8 {
			return
		}
		if _, err := strconv.Atoi(fields[0]); err != nil {
			return
		}
		ts := labels{"name": replace(fields[2])}
		add("chassis_volts", severity(fields[1]), ts, descDellHWVolt)
		if i, err := extract(fields[3], "V"); err == nil {
			add("chassis_volts_reading", i, ts, descDellHWVoltReadings)
		}
	}, "chassis", "volts")
	return nil
}
