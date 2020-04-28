package main

import (
	"fmt"
	"github.com/foreso-GitHub/jingtum-monitor/demo/collector"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

func init() {
	//注册自身采集器
	//prometheus.MustRegister(collector.NewNodeCollector())
	prometheus.MustRegister(collector.NewJtCollector())
}

func main() {
	collector.GetBlockNumber()
	http.Handle("/metrics", promhttp.Handler())
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Error occur when start server %v", err)
	}

}
