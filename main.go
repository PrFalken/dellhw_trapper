package main

import (
	"flag"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	listenAddress     = flag.String("web.listen", ":4242", "Address on which to expose metrics and web interface.")
	metricsPath       = flag.String("web.path", "/metrics", "Path under which to expose metrics.")
	enabledCollectors = flag.String("collect", "dummy,chassis,memory,processors,ps,ps_amps_sysboard_pwr,storage_battery,storage_enclosure,storage_vdisk,system,temps,volts", "Comma-separated list of collectors to use.")

	collectors = map[string]Collector{
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
)

type Collector struct {
	F func() error
}

func Add(name string, value string, t prometheus.Labels, desc string) {
	log.Println("Adding metric : ", name, t, value)
	d := prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace:   "dell",
		Subsystem:   "hw",
		Name:        name,
		Help:        desc,
		ConstLabels: t,
	})
	floatValue, err := strconv.ParseFloat(value, 64)
	if err != nil {
		log.Println("Could not parse value for metric ", name)
		return
	}
	d.Set(floatValue)
	prometheus.MustRegister(d)
}

func collect(collectors map[string]Collector) {
	for _, name := range strings.Split(*enabledCollectors, ",") {
		collector := collectors[name]
		log.Println("Running collector ", name)
		err := collector.F()
		if err != nil {
			log.Println(err)
		}
	}
}

func main() {
	flag.Parse()

	collect(collectors)
	http.Handle(*metricsPath, prometheus.Handler())

	log.Print("listening to ", *listenAddress)
	log.Fatal(http.ListenAndServe(*listenAddress, nil))

}
