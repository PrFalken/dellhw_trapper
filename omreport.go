package main

import (
	"encoding/json"
	"strconv"
	"strings"

	"bosun.org/metadata"
	"bosun.org/slog"
	"bosun.org/util"
)

func readOmreport(f func([]string), args ...string) {
	args = append(args, "-fmt", "ssv")
	_ = util.ReadCommand(func(line string) error {
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
	Add(&md, "hw.dummy", "0", TagSet{"test": "dummy"}, metadata.Gauge, metadata.Ok, "Dummy description")
	return md, nil
}

func c_omreport_chassis() (MultiDataPoint, error) {
	var md MultiDataPoint
	readOmreport(func(fields []string) {
		if len(fields) != 2 || fields[0] == "SEVERITY" {
			return
		}
		component := strings.Replace(fields[1], " ", "_", -1)
		Add(&md, "hw.chassis", severity(fields[0]), TagSet{"component": component}, metadata.Gauge, metadata.Ok, descDellHWChassis)
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
		Add(&md, "hw.system", severity(fields[0]), TagSet{"component": component}, metadata.Gauge, metadata.Ok, descDellHWSystem)
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
		Add(&md, "hw.storage.enclosure", fields[1], TagSet{"id": id}, metadata.Gauge, metadata.Ok, descDellHWStorageEnc)
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
		Add(&md, "hw.storage.vdisk", fields[1], TagSet{"id": id}, metadata.Gauge, metadata.Ok, descDellHWVDisk)
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
		Add(&md, "hw.ps", fields[1], ts, metadata.Gauge, metadata.Ok, descDellHWPS)
		pm := &metadata.HWPowerSupplyMeta{}
		if len(fields) < 6 {
			return
		}
		if fields[4] != "" {
			pm.RatedInputWattage = fields[4]
		}
		if fields[5] != "" {
			pm.RatedOutputWattage = fields[5]
		}
		if j, err := json.Marshal(&pm); err == nil {
			AddMeta("", ts, "psMeta", string(j), true)
		} else {
			slog.Error(err)
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
			Add(&md, "hw.chassis.current.reading", v_fields[0], TagSet{"id": id}, metadata.Gauge, metadata.A, descDellHWCurrent)
		} else if len(fields) == 6 && (fields[2] == "System Board Pwr Consumption" || fields[2] == "System Board System Level") {
			v_fields := strings.Fields(fields[3])
			warn_fields := strings.Fields(fields[4])
			fail_fields := strings.Fields(fields[5])
			if len(v_fields) < 2 || len(warn_fields) < 2 || len(fail_fields) < 2 {
				return
			}
			Add(&md, "hw.chassis.power.reading", v_fields[0], nil, metadata.Gauge, metadata.Watt, descDellHWPower)
			Add(&md, "hw.chassis.power.warn_level", warn_fields[0], nil, metadata.Gauge, metadata.Watt, descDellHWPowerThreshold)
			Add(&md, "hw.chassis.power.fail_level", fail_fields[0], nil, metadata.Gauge, metadata.Watt, descDellHWPowerThreshold)
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
		Add(&md, "hw.storage.battery", fields[1], TagSet{"id": id}, metadata.Gauge, metadata.Ok, descDellHWStorageBattery)
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
		Add(&md, "hw.storage.controller", fields[1], ts, metadata.Gauge, metadata.Ok, descDellHWStorageCtl)
		cm := &metadata.HWControllerMeta{}
		if len(fields) < 8 {
			return
		}
		if fields[2] != "" {
			cm.Name = fields[2]
		}
		if fields[3] != "" {
			cm.SlotId = fields[3]
		}
		if fields[4] != "" {
			cm.State = fields[4]
		}
		if fields[5] != "" {
			cm.FirmwareVersion = fields[5]
		}
		if fields[7] != "" {
			cm.DriverVersion = fields[7]
		}
		if j, err := json.Marshal(&cm); err == nil {
			AddMeta("", ts, "controllerMeta", string(j), true)
		} else {
			slog.Error(err)
		}
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
		Add(md, "hw.storage.pdisk", fields[1], ts, metadata.Gauge, metadata.Ok, descDellHWPDisk)
		if len(fields) < 32 {
			return
		}
		dm := &metadata.HWDiskMeta{}
		if fields[2] != "" {
			dm.Name = fields[2]
		}
		if fields[6] != "" {
			dm.Media = fields[6]
		}
		if fields[19] != "" {
			dm.Capacity = fields[19]
		}
		if fields[23] != "" {
			dm.VendorId = fields[23]
		}
		if fields[24] != "" {
			dm.ProductId = fields[24]
		}
		if fields[25] != "" {
			dm.Serial = fields[25]
		}
		if fields[26] != "" {
			dm.Part = fields[26]
		}
		if fields[27] != "" {
			dm.NegotatiedSpeed = fields[27]
		}
		if fields[28] != "" {
			dm.CapableSpeed = fields[28]
		}
		if fields[31] != "" {
			dm.SectorSize = fields[31]

		}
		if j, err := json.Marshal(&dm); err == nil {
			AddMeta("", ts, "physicalDiskMeta", string(j), true)
		} else {
			slog.Error(err)
		}
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
		Add(&md, "hw.chassis.processor", fields[1], ts, metadata.Gauge, metadata.Ok, descDellHWCPU)
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
		Add(&md, "hw.chassis.fan", fields[1], ts, metadata.Gauge, metadata.Ok, descDellHWFan)
		fs := strings.Fields(fields[3])
		if len(fs) == 2 && fs[1] == "RPM" {
			i, err := strconv.Atoi(fs[0])
			if err == nil {
				Add(&md, "hw.chassis.fan.reading", i, ts, metadata.Gauge, metadata.RPM, descDellHWFanSpeed)
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
		Add(&md, "hw.chassis.memory", fields[1], ts, metadata.Gauge, metadata.Ok, descDellHWMemory)
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
		Add(&md, "hw.chassis.temps", fields[1], ts, metadata.Gauge, metadata.Ok, descDellHWTemp)
		fs := strings.Fields(fields[3])
		if len(fs) == 2 && fs[1] == "C" {
			i, err := strconv.ParseFloat(fs[0], 64)
			if err == nil {
				Add(&md, "hw.chassis.temps.reading", i, ts, metadata.Gauge, metadata.C, descDellHWTempReadings)
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
		Add(&md, "hw.chassis.volts", fields[1], ts, metadata.Gauge, metadata.Ok, descDellHWVolt)
		if i, err := extract(fields[3], "V"); err == nil {
			Add(&md, "hw.chassis.volts.reading", i, ts, metadata.Gauge, metadata.V, descDellHWVoltReadings)
		}
	}, "chassis", "volts")
	return md, nil
}
