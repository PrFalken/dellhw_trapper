package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"bosun.org/metadata"
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
	listenAddress     = flag.String("web.listen", ":4242", "Address on which to expose metrics and web interface.")
	metricsPath       = flag.String("web.path", "/metrics", "Path under which to expose metrics.")
	enabledCollectors = flag.String("collect", "dummy,chassis,memory,processors", "Comma-separated list of collectors to use.")
)

type MetricStorage struct {
	collections map[string]*MultiDataPoint
	lock        sync.RWMutex
}

func NewMetricStorage() *MetricStorage {
	var m MetricStorage
	m.collections = make(map[string]*MultiDataPoint)
	return &m
}

type DataPoint struct {
	Metric    string      `json:"metric"`
	Timestamp int64       `json:"timestamp"`
	Value     interface{} `json:"value"`
	Tags      TagSet      `json:"tags"`
}

type MultiDataPoint []*DataPoint

type TagSet map[string]string

type Collector struct {
	F func() (MultiDataPoint, error)
}

func Add(md *MultiDataPoint, name string, value interface{}, t TagSet, rate metadata.RateType, unit metadata.Unit, desc string) {
	log.Println("Adding metric : ", name, t, value)

	d := DataPoint{
		Metric: name,
		Value:  value,
		Tags:   t,
	}
	*md = append(*md, &d)
}

func AddMeta(metric string, tags TagSet, name string, value interface{}, setHost bool) {
	// Keeping it for future use
	// fmt.Println(metric, tags, name, value, setHost)
}

// extract tries to return a parsed number from s with given suffix. A space may
// be present between number ond suffix.
func extract(s, suffix string) (float64, error) {
	if !strings.HasSuffix(s, suffix) {
		return 0, fmt.Errorf("extract: suffix not found")
	}
	s = s[:len(s)-len(suffix)]
	return strconv.ParseFloat(strings.TrimSpace(s), 64)
}

// severity returns 0 if s is not "Ok" or "Non-Critical", else 1.
func severity(s string) int {
	if s != "Ok" && s != "Non-Critical" {
		return 1
	}
	return 0
}

func (s *MetricStorage) collect(collectors map[string]Collector) {
	for {
		for _, name := range strings.Split(*enabledCollectors, ",") {
			collector := collectors[name]
			log.Println("Running collector ", name)
			md, err := collector.F()
			if err != nil {
				log.Println(err)
			}
			s.lock.Lock()
			s.collections[name] = &md
			s.lock.Unlock()
		}
		time.Sleep(60 * time.Second)
	}

}

func main() {
	flag.Parse()

	storage := NewMetricStorage()
	collectors := map[string]Collector{
		"dummy":      Collector{F: dummy_report},
		"chassis":    Collector{F: c_omreport_chassis},
		"fans":       Collector{F: c_omreport_fans},
		"memory":     Collector{F: c_omreport_memory},
		"processors": Collector{F: c_omreport_processors},
		"ps":         Collector{F: c_omreport_ps},
		"ps_amps_sysboard_pwr": Collector{F: c_omreport_ps_amps_sysboard_pwr},
		"storage_battery":      Collector{F: c_omreport_storage_battery},
		"storage_controller":   Collector{F: c_omreport_storage_controller},
		"storage_enclosure":    Collector{F: c_omreport_storage_enclosure},
		"storage_vdisk":        Collector{F: c_omreport_storage_vdisk},
		"system":               Collector{F: c_omreport_system},
		"temps":                Collector{F: c_omreport_temps},
		"volts":                Collector{F: c_omreport_volts},
	}

	go storage.collect(collectors)

	// SERVE
	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		for _, md := range storage.collections {
			for _, d := range *md {
				tags, err := json.Marshal(d.Tags)
				if err != nil {
					http.Error(w, err.Error(), 500)
					return
				}
				fmt.Fprint(w, d.Metric, string(tags), " ", d.Value)
			}
		}
	})

	log.Print("listening to ", *listenAddress)
	log.Fatal(http.ListenAndServe(*listenAddress, nil))

}
