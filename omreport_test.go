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
	// Fake "omreport storage enclosure" splitted string return"
	if reflect.DeepEqual(args, []string{"storage", "enclosure"}) {
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
