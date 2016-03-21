package main

import (
	"log"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

func addToPrometheus(name string, value string, t prometheus.Labels, desc string) {
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
