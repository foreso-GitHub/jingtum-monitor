package main

import (
	"fmt"
	"github.com/foreso-GitHub/jingtum-monitor/common"
	"github.com/foreso-GitHub/jingtum-monitor/exporter"
	"time"
)

var firstCollect = true
var tpsStatus = exporter.CreateJtTpsStatus(-1)

func test() {
	config := common.LoadConfig("./config/config.json")

	//region test
	testRun(config)
	time.Sleep(6000000000)
	testRun(config)
	time.Sleep(61000000000)
	testRun(config)
	//endregion
}

func testRun(config common.ExporterConfig) {
	//url := "http://" + "box-admin.elerp.net" + ":" + "10201" + "/v1/jsonrpc" //请求地址
	//url := "http://" + "180.76.125.22" + ":" + "9545" + "/v1/jsonrpc" //请求地址
	//testLibrary(url)
	//testTps(url)

	network := exporter.LoadJtNetworkConfig(config.JtConfigPath)
	nodes := network.NodeList
	url := exporter.GetRandUrl(nodes)

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
	if firstCollect {
		blockNumber, _ := exporter.GetBlockNumber(url)
		tpsStatus = exporter.CreateJtTpsStatus(blockNumber)
		firstCollect = false
	}

	flushOK := exporter.FlushTpsStatus(url, tpsStatus)
	fmt.Printf("flushOK: %+v\n", flushOK)
	fmt.Printf("===CurrentBlockNumber: %+v\n", tpsStatus.CurrentBlockNumber)
}
