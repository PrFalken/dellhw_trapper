package main

import (
	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	// RootCmd is the main command
	RootCmd = &cobra.Command{
		Use:   "dellhw_trapper",
		Short: "Zabbix exporter for Dell Hardware components",
		Run: func(cmd *cobra.Command, args []string) {
			runMainCommand()
		},
	}

	HWEVersion string
	BuildDate  string
	logLevel   string

	enabledCollectors   string
	discoveryNameSpace  string
	zabbixFromHost      string
	zabbixServerAddress string
	zabbixServerPort    string
	zabbixDiscovery     bool
	zabbixUpdateItems   bool

	cache          = newMetricStorage()
	metricCounts   = make(map[string]int)
	metricStatuses = make(map[string]int)
)

func init() {
	RootCmd.Flags().StringVarP(&logLevel, "loglevel", "L", "info", "Set log level")
	RootCmd.Flags().StringVarP(&enabledCollectors, "collect", "c", "chassis,fans,memory,processors,ps,ps_amps_sysboard_pwr,storage_battery,storage_enclosure,storage_controller,storage_vdisk,system,temps,volts", "Comma-separated list of collectors to use.")
	RootCmd.Flags().StringVarP(&zabbixFromHost, "zabbix-from", "f", getFQDN(), "Send to Zabbix from this host name. You can also set HOSTNAME and DOMAINNAME environment variables.")
	RootCmd.Flags().StringVarP(&zabbixServerAddress, "zabbix-server", "z", "localhost", "Zabbix server hostname or address")
	RootCmd.Flags().StringVarP(&zabbixServerPort, "zabbix-port", "p", "10051", "Zabbix server port")
	RootCmd.Flags().StringVarP(&discoveryNameSpace, "namespace", "n", "", "Discovery key")
	RootCmd.Flags().BoolVar(&zabbixDiscovery, "discovery", false, "Perform Zabbix low level discovery on hardware elements")
	RootCmd.Flags().BoolVar(&zabbixUpdateItems, "update-items", false, "Get & send items to Zabbix. This is the default behaviour")
	RootCmd.AddCommand(versionCmd)

}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of hardware_exporter",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print("hardware_exporter\n\n")
		fmt.Printf("version    : %s\n", HWEVersion)
		if BuildDate != "" {
			fmt.Printf("build date : %s\n", BuildDate)
		}
	},
}

func runMainCommand() {

	if logLevel == "info" {
		log.SetLevel(log.InfoLevel)
	}
	if logLevel == "debug" {
		log.SetLevel(log.DebugLevel)
	}
	if logLevel == "error" {
		log.SetLevel(log.ErrorLevel)
	}

	err := collect(collectors)
	if err != nil {
		log.Debug("Collect failed")
		os.Exit(1)
	}

	if zabbixDiscovery {
		discovery()
	} else {
		updateItems()
	}
}

func main() {

	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

}
