package main

import (
	"fmt"
	"github.com/foreso-GitHub/jingtum-monitor/demo/demo1/collector"
	"github.com/foreso-GitHub/jingtum-monitor/exporter"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

func init() {
	//注册自身采集器
	prometheus.MustRegister(collector.NewNodeCollector())
	prometheus.MustRegister(exporter.NewJtCollector())
}

func main() {
	http.Handle("/metrics", promhttp.Handler())
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Error occur when start server %v", err)
	}
}
