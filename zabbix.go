package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"

	zabbix "github.com/AlekSi/zabbix-sender"
	log "github.com/Sirupsen/logrus"
)

type zabbixDiscoveryItem struct {
	Name string `json:"{#DELLHWCOMPONENTNAME}"`
}

func newZabbixDiscoveryItem(name string) *zabbixDiscoveryItem {
	item := zabbixDiscoveryItem{
		Name: name,
	}
	return &item
}

func sendToZabbix() {
	cache.Lock.Lock()
	initValue := make(map[string]interface{})
	di := zabbix.MakeDataItems(initValue, zabbixFromHost)

	metricPrefix := "dellhw.components"

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
			log.Debug("Discovery failure, could not marshal to json")
			fmt.Println("2")
			os.Exit(2)
		}

		discoveryPayload := make(map[string]interface{})
		discoveryPayload[metricPrefix+".discovery"] = string(jsonOutput)

		di = zabbix.MakeDataItems(discoveryPayload, zabbixFromHost)

	} else {

		// add discovery name wrap
		newMap := make(map[string]interface{})
		for k, v := range cache.metrics {
			newKey := metricPrefix + "[" + k + "]"
			newMap[newKey] = v
		}

		di = zabbix.MakeDataItems(newMap, zabbixFromHost)
	}
	cache.Lock.Unlock()
	addr, _ := net.ResolveTCPAddr("tcp", zabbixServerAddress+":"+zabbixServerPort)
	res, err := zabbix.Send(addr, di)
	if err != nil {
		log.Debug("Step 4 - Sent to Zabbix Server failed : ", err)
		fmt.Println("4")
		os.Exit(4)
	}
	log.Debug(*res)
	fmt.Println("0")
}
