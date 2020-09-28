package main

import (
	"fmt"
	"github.com/foreso-GitHub/jingtum-monitor/common"
	"github.com/foreso-GitHub/jingtum-monitor/exporter"
	"time"
)

var firstCollect = true
var tpsStatus = exporter.CreateJtTpsStatus(-1)
var connectedUrl = ""
var config = common.LoadConfig("./config/config.json")

func test() {
	config := common.LoadConfig("./config/config.json")

	//region test
	for i := 0; i < 100; i++ {
		testRunTps(config)
		time.Sleep(6000000000)
	}

	//endregion
}

func testRunTps(config common.ExporterConfig) {
	//url := "http://" + "box-admin.elerp.net" + ":" + "10201" + "/v1/jsonrpc" //请求地址
	//url := "http://" + "180.76.125.22" + ":" + "9545" + "/v1/jsonrpc" //请求地址
	//testLibrary(url)
	//testTps(url)

	//network := exporter.LoadJtNetworkConfig(config.JtConfigPath)
	//nodes := network.NodeList
	//url := exporter.GetRandUrl(nodes)

	url, _, _ := exporter.GetBlockNumberByRandNode()
	//fmt.Printf("blockNumber: %+v\n", blockNumber)
	testTps(url)
	//testLibrary(url)
}

func testLibrary(url string) {
	fmt.Printf("Url: %v", url)

	blockNumber, err := exporter.GetBlockNumberByNode(url)
	fmt.Printf("blockNumber: %v\n", blockNumber)
	fmt.Printf("blockNumber err: %v\n", err)
	_, block, err := exporter.GetBlockByNumberByRandNode(blockNumber)
	fmt.Printf("block: %+v\n", block)
	fmt.Printf("block err: %v\n", err)
}

func testTps(url string) {
	fmt.Printf("Url: %v\n", url)
	if firstCollect {
		blockNumber, _ := exporter.GetBlockNumberByNode(url)
		//fmt.Printf("blockNumber: %+v\n", blockNumber)
		tpsStatus = exporter.CreateJtTpsStatus(blockNumber)
		firstCollect = false
	}

	_ = exporter.FlushTpsStatus(tpsStatus)
	//fmt.Printf("flushOK: %+v\n", flushOK)
	//fmt.Printf("TpsMap: %+v\n", tpsStatus.TpsMap)
	//fmt.Printf("===CurrentBlockNumber: %+v\n", tpsStatus.CurrentBlockNumber)
}
