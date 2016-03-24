package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"

	zabbix "github.com/AlekSi/zabbix-sender"
	log "github.com/Sirupsen/logrus"
)

type zabbixItem struct {
	Name        string
	Labels      map[string]string
	Value       interface{}
	Description string
}

func newZabbixItem(name string, labels labels, value interface{}, desc string) *zabbixItem {
	item := zabbixItem{
		Name:        name,
		Labels:      labels,
		Value:       value,
		Description: desc,
	}
	return &item
}

func discovery() {
	log.Debug("Running discovery")
	cache.Lock.Lock()
	defer cache.Lock.Unlock()
	discoData := make(map[string][]labels)
	discoItemList := []labels{}
	for _, item := range cache.metrics {
		discoItem := item.Labels
		discoItemList = append(discoItemList, discoItem)
	}
	discoData["data"] = discoItemList

	jsonOutput, err := json.Marshal(discoData)
	if err != nil {
		log.Debug("Discovery failure, could not marshal to json")
		fmt.Println("2")
		os.Exit(2)
	}

	discoveryPayload := make(map[string]interface{})
	discoveryPayload[discoveryNameSpace+".discovery"] = string(jsonOutput)
	log.Debug(discoveryPayload)
	di := zabbix.MakeDataItems(discoveryPayload, zabbixFromHost)
	sendToZabbix(di)
}

func updateItems() {
	log.Debug("Running update-items")
	cache.Lock.Lock()
	defer cache.Lock.Unlock()

	// add discovery name wrap
	newMap := make(map[string]interface{})
	for _, metric := range cache.metrics {
		key := metric.Name
		newMap[key] = metric.Value
	}
	log.Debug("sending items : ", newMap)
	di := zabbix.MakeDataItems(newMap, zabbixFromHost)
	sendToZabbix(di)
}

func sendToZabbix(di zabbix.DataItems) {
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
