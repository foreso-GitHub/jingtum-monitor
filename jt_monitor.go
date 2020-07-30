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
	//prometheus.MustRegister(collector.NewNodeCollector())
	prometheus.MustRegister(exporter.NewJtCollector())
}

func main() {

	config := common.LoadConfig("./config/config.json")

	//region test
	test()
	//network := exporter.LoadJtNetworkConfig(config.JtConfigPath)
	//exporter.FlushNetwork(&network)
	//endregion

	http.Handle("/metrics", promhttp.Handler())
	if err := http.ListenAndServe(config.ExportAddress, nil); err != nil {
		fmt.Printf("Error occur when start server %v", err)
	}
}

func test() {
	//url := "http://" + "box-admin.elerp.net" + ":" + "10201" + "/v1/jsonrpc" //请求地址
	url := "http://" + "180.76.125.22" + ":" + "9545" + "/v1/jsonrpc" //请求地址
	testLibrary(url)
	testTps(url)
}

func testLibrary(url string) {
	fmt.Printf("Url: %v", url)

	blockNumber, err := exporter.GetBlockNumber(url)
	fmt.Printf("blockNumber: %v\n", blockNumber)
	fmt.Printf("blockNumber err: %v\n", err)
	block, err := exporter.GetBlockByNumber(url, blockNumber)
	fmt.Printf("block: %+v\n", block)
	fmt.Printf("block err: %v\n", err)
}

func testTps(url string) {
	fmt.Printf("Url: %v\n", url)
	status := exporter.InitJtTpsStatus()
	flushOK := exporter.FlushTpsStatus(url, status)
	fmt.Printf("flushOK: %+v\n", flushOK)
	fmt.Printf("flushOK: %+v\n", status)
}
