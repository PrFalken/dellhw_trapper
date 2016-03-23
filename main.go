package main

import (
	"fmt"
	"net/http"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/cobra"
)

var (
	// RootCmd is the main command
	RootCmd = &cobra.Command{
		Use:   "hardware_exporter",
		Short: "Prometheus and Zabbix exporter for Dell Hardware components",
		Run: func(cmd *cobra.Command, args []string) {
			runMainCommand()
		},
	}

	exporterType        string
	listenAddress       string
	metricsPath         string
	enabledCollectors   string
	zabbixFromHost      string
	zabbixServerAddress string
	zabbixServerPort    string
	zabbixDiscovery     bool
	zabbixUpdateItems   bool

	cache = newMetricStorage()
)

func init() {
	RootCmd.Flags().StringVarP(&exporterType, "type", "t", "zabbix", "Exporter type : prometheus or zabbix")
	RootCmd.Flags().StringVarP(&listenAddress, "web-listen", "l", "127.0.0.1", "Address on which to expose metrics and web interface.")
	RootCmd.Flags().StringVarP(&metricsPath, "web-path", "m", "/metrics", "Path under which to expose metrics.")
	RootCmd.Flags().StringVarP(&enabledCollectors, "collect", "c", "chassis,fans,memory,processors,ps,ps_amps_sysboard_pwr,storage_battery,storage_enclosure,storage_vdisk,system,temps,volts", "Comma-separated list of collectors to use.")
	RootCmd.Flags().StringVarP(&zabbixFromHost, "zabbix-from", "f", getFQDN(), "Send to Zabbix from this host name. You can also set HOSTNAME and DOMAINNAME environment variables.")
	RootCmd.Flags().StringVarP(&zabbixServerAddress, "zabbix-server", "z", "localhost", "Zabbix server hostname or address")
	RootCmd.Flags().StringVarP(&zabbixServerPort, "zabbix-port", "p", "10051", "Zabbix server port")
	RootCmd.Flags().BoolVar(&zabbixDiscovery, "discovery", false, "Perform Zabbix low level discovery on hardware elements")
	RootCmd.Flags().BoolVar(&zabbixUpdateItems, "update-items", false, "Get & send items to Zabbix. This is the default behaviour")

	log.SetLevel(log.DebugLevel)

}

func runMainCommand() {
	err := collect(collectors)
	if err != nil {
		log.Info("Collect failed")
		os.Exit(1)
	}

	switch exporterType {
	case "prometheus":
		http.Handle(metricsPath, prometheus.Handler())
		log.Info("listening to ", listenAddress)
		log.Fatal(http.ListenAndServe(listenAddress, nil))

	case "zabbix":
		sendToZabbix()
	}
}

func main() {

	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

}
