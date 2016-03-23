package main

import (
	"encoding/json"
	"net"
	"os"
	"strings"

	zabbix "github.com/AlekSi/zabbix-sender"
	log "github.com/Sirupsen/logrus"
	"github.com/prometheus/client_golang/prometheus"
)

type zabbixDiscoveryItem struct {
	Name string `json:"#DELLHWCOMPONENTNAME"`
}

func newZabbixDiscoveryItem(name string) *zabbixDiscoveryItem {
	item := zabbixDiscoveryItem{
		Name: name,
	}
	return &item
}

func addToZabbix(name string, value string, t prometheus.Labels) {
	cache.Lock.Lock()
	defer cache.Lock.Unlock()
	zabbixMetricName := "hw." + strings.Replace(name, "_", ".", -1)
	for _, v := range t {
		zabbixMetricName += "." + v
	}
	cache.metrics[zabbixMetricName] = value
}

func sendToZabbix() {
	cache.Lock.Lock()
	initValue := make(map[string]interface{})
	di := zabbix.MakeDataItems(initValue, zabbixFromHost)

	if zabbixDiscovery {
		discoData := make(map[string][]zabbixDiscoveryItem)
		discoItemList := []zabbixDiscoveryItem{}
		for metricName := range cache.metrics {
			discoItem := newZabbixDiscoveryItem(metricName)
			discoItemList = append(discoItemList, *discoItem)
		}
		discoData["data"] = discoItemList

		jsonOutput, err := json.Marshal(discoData)
		if err != nil {
			log.Error("Discovery failure, could not marshal to json")
			os.Exit(2)
		}

		discoveryPayload := make(map[string]interface{})
		discoveryPayload["dellhw.components.discovery"] = string(jsonOutput)

		di = zabbix.MakeDataItems(discoveryPayload, zabbixFromHost)

	} else {
		di = zabbix.MakeDataItems(cache.metrics, zabbixFromHost)
	}
	cache.Lock.Unlock()
	addr, _ := net.ResolveTCPAddr("tcp", zabbixServerAddress+":"+zabbixServerPort)
	res, err := zabbix.Send(addr, di)
	if err != nil {
		log.Error("Step 4 - Sent to Zabbix Server failed : ", err)
		os.Exit(4)
	}
	log.Info(*res)
}
