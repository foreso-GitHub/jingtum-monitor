package main

import (
	"fmt"
	"github.com/foreso-GitHub/jingtum-monitor/common"
	"github.com/foreso-GitHub/jingtum-monitor/exporter"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

func init() {
	//注册自身采集器
	prometheus.MustRegister(exporter.NewJtCollector())
}

func main() {
	//debug()

	config := common.LoadConfig("./config/config.json")
	http.Handle("/metrics", promhttp.Handler())
	if err := http.ListenAndServe(config.ExportAddress, nil); err != nil {
		fmt.Printf("Error occur when start server %v", err)
	}
}
