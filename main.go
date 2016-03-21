package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	exporterType        = flag.String("type", "prometheus", "Exporter type : prometheus or zabbix")
	listenAddress       = flag.String("web.listen", ":4242", "Address on which to expose metrics and web interface.")
	metricsPath         = flag.String("web.path", "/metrics", "Path under which to expose metrics.")
	enabledCollectors   = flag.String("collect", "dummy,chassis,memory,processors,ps,ps_amps_sysboard_pwr,storage_battery,storage_enclosure,storage_vdisk,system,temps,volts", "Comma-separated list of collectors to use.")
	zabbixFromHost      = flag.String("zabbix.from", getFQDN(), "Send to Zabbix from this host name. You can also set HOSTNAME and DOMAINNAME environment variables.")
	zabbixServerAddress = flag.String("zabbix.server.address", "localhost", "Zabbix server hostname or address")
	zabbixServerPort    = flag.String("zabbix.server.port", "10051", "Zabbix server port")
	zabbixDiscovery     = flag.Bool("zabbix.discovery", false, "Perform Zabbix low level discovery on hardware elements")
	cache               = newMetricStorage()
)

type metricStorage struct {
	Lock    sync.RWMutex
	metrics map[string]interface{}
}

func newMetricStorage() *metricStorage {
	ms := new(metricStorage)
	ms.metrics = make(map[string]interface{})
	return ms
}

func add(name string, value string, t prometheus.Labels, desc string) {

	switch *exporterType {

	case "prometheus":
		addToPrometheus(name, value, t, desc)

	case "zabbix":
		addToZabbix(name, value, t)
	}
}

func main() {
	flag.Parse()
	err := collect(collectors)
	if err != nil {
		log.Println("Collect failed")
		os.Exit(1)
	}

	switch *exporterType {
	case "prometheus":
		http.Handle(*metricsPath, prometheus.Handler())
		log.Println("listening to ", *listenAddress)
		log.Fatal(http.ListenAndServe(*listenAddress, nil))

	case "zabbix":
		sendToZabbix()
	}
}
