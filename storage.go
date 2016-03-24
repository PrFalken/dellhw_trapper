package main

import "sync"

type metricStorage struct {
	Lock    sync.RWMutex
	metrics map[string]zabbixItem
}

func newMetricStorage() *metricStorage {
	ms := new(metricStorage)
	ms.metrics = make(map[string]zabbixItem)
	return ms
}
