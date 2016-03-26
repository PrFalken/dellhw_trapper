package main

import (
	"fmt"
	"reflect"
	"testing"
)

func showCache() {
	for _, m := range cache.metrics {
		fmt.Printf("%+v\n", m)
	}
}

type testOmReport struct{}

func newTestOmReport() *testOmReport {
	return &testOmReport{}
}

//
// Mock omreport command

func (o *testOmReport) Report(f func([]string), args ...string) {
	// Fake "omreport chassis" splitted string return
	if reflect.DeepEqual(args, []string{"chassis"}) {
		sp := []string{"Ok", "testChassisName"}
		f(sp)
	}
	// Fake "omreport system" splitted string return
	if reflect.DeepEqual(args, []string{"system"}) {
		sp := []string{"Ok", "testSystemName"}
		f(sp)
	}
	// Fake "omreport storage enclosure" splitted string return
	if reflect.DeepEqual(args, []string{"storage", "enclosure"}) {
		sp := []string{"0:1", "Ok", "blah"}
		f(sp)
	}
	// Fake "omreport storage vdisk" splitted string return
	if reflect.DeepEqual(args, []string{"storage", "vdisk"}) {
		sp := []string{"0:1", "Ok", "blah"}
		f(sp)
	}
	// Fake "omreport chassis pwrsupplies" splitted string return
	if reflect.DeepEqual(args, []string{"chassis", "pwrsupplies"}) {
		sp := []string{"0:1", "Ok", "blah", "bloh", "42 W", "4242 W"}
		f(sp)
	}
	// Fake "omreport chassis pwrmonitoring splitted string return
	if reflect.DeepEqual(args, []string{"chassis", "pwrmonitoring"}) {
		sp := []string{"0:1", "Ok", "blah"}
		f(sp)
	}

}

//
// Test for OmReport functions

func TestOmReportChassis(t *testing.T) {
	to := newTestOmReport()
	omreportChassis(to)
	returnedLabel := cache.metrics["chassis"].Labels["component"]
	if returnedLabel != "testChassisName" {
		t.Error("Expected testChassisName, got ", returnedLabel)
	}
	value := cache.metrics["chassis"].Value
	if value != "0" {
		t.Error("Expected return value 0, got ", value)
	}
}

func TestOmReportSystem(t *testing.T) {
	to := newTestOmReport()
	omreportSystem(to)
	returnedLabel := cache.metrics["system"].Labels["component"]
	if returnedLabel != "testSystemName" {
		t.Error("Expected testSystemName, got ", returnedLabel)
	}
	returnedValue := cache.metrics["system"].Value
	if returnedValue != "0" {
		t.Error("Expected return value 0, got ", returnedValue)
	}
}

func TestOmreportStorageEnclosure(t *testing.T) {
	to := newTestOmReport()
	omreportStorageEnclosure(to)
	returnedLabel := cache.metrics["storage_enclosure"].Labels["id"]
	if returnedLabel != "0_1" {
		t.Error("Expected 0_1, got ", returnedLabel)
	}
	returnedValue := cache.metrics["storage_enclosure"].Value
	if returnedValue != "0" {
		t.Error("Expected return value 0, got ", returnedValue)
	}

}

func TestOmreportStorageVdisk(t *testing.T) {
	to := newTestOmReport()
	omreportStorageVdisk(to)
	returnedLabel := cache.metrics["dell.hardware.raid.logicaldrive[0_1,status]"].Labels["{#LOGICALDRIVESLOT}"]
	if returnedLabel != "0_1" {
		t.Error("Expected 0_1, got ", returnedLabel)
	}
	returnedValue := cache.metrics["dell.hardware.raid.logicaldrive[0_1,status]"].Value
	if returnedValue != "0" {
		t.Error("Expected return value 0, got ", returnedValue)
	}

}

func TestOmreportPs(t *testing.T) {
	to := newTestOmReport()
	omreportPs(to)
	returnedLabel := cache.metrics["dell.hardware.power[0_1,status]"].Labels["{#POWERSLOT}"]
	if returnedLabel != "0_1" {
		t.Error("Expected 0_1, got ", returnedLabel)
	}
	returnedStatusValue := cache.metrics["dell.hardware.power[0_1,status]"].Value
	if returnedStatusValue != "0" {
		t.Error("Expected return value 0, got ", returnedStatusValue)
	}
	returnediWattsValue := cache.metrics["dell.hardware.power[0_1,input_watts]"].Value
	if returnediWattsValue != "42" {
		t.Error("Expected return value 42, got ", returnediWattsValue)
	}
	returnedoWattsValue := cache.metrics["dell.hardware.power[0_1,output_watts]"].Value
	if returnedoWattsValue != "4242" {
		t.Error("Expected return value 4242, got ", returnedoWattsValue)
	}
}
