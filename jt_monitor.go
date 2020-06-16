package main

import (
	"fmt"
	"github.com/foreso-GitHub/jingtum-monitor/common"
	"github.com/foreso-GitHub/jingtum-monitor/exporter"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

//import (
//	"fmt"
//	"github.com/foreso-GitHub/jingtum-monitor/common"
//	"github.com/foreso-GitHub/jingtum-monitor/exporter"
//	"github.com/prometheus/client_golang/prometheus"
//	"github.com/prometheus/client_golang/prometheus/promhttp"
//	"net/http"
//)

//import (
//	"fmt"
//	"github.com/foreso-GitHub/jingtum-monitor/common"
//	"github.com/foreso-GitHub/jingtum-monitor/exporter"
//	"github.com/prometheus/client_golang/prometheus"
//	"github.com/prometheus/client_golang/prometheus/promhttp"
//	"net/http"
//)

func init() {
	//注册自身采集器
	//prometheus.MustRegister(collector.NewNodeCollector())
	prometheus.MustRegister(exporter.NewJtCollector())
}

func main() {
	test()

	config := common.LoadConfig("./config/config.json")
	http.Handle("/metrics", promhttp.Handler())
	if err := http.ListenAndServe(config.ExportAddress, nil); err != nil {
		fmt.Printf("Error occur when start server %v", err)
	}
}

func test() {
	testLibrary()
	testTps()
}

func testLibrary() {
	url := "http://" + "box-admin.elerp.net" + ":" + "10201" + "/v1/jsonrpc" //请求地址
	fmt.Printf("Url: %v", url)

	blockNumber, err := exporter.GetBlockNumber(url)
	fmt.Println("blockNumber: %v", blockNumber)
	fmt.Println("blockNumber err: %v", err)
	block, err := exporter.GetBlockByNumber(url, blockNumber)
	fmt.Println("block: %+v", block)
	fmt.Println("block err: %v", err)
}

func testTps() {
	url := "http://" + "box-admin.elerp.net" + ":" + "10201" + "/v1/jsonrpc" //请求地址
	fmt.Println("Url: %v", url)
	status := exporter.InitJtTpsStatus()
	flushOK := exporter.FlushTpsStatus(url, status)
	fmt.Println("flushOK: %+v", flushOK)
	fmt.Println("flushOK: %+v", status)
}
