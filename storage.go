package main

type metricStorage struct {
	metrics map[string]zabbixItem
}

func newMetricStorage() *metricStorage {
	ms := new(metricStorage)
	ms.metrics = make(map[string]zabbixItem)
	return ms
}
