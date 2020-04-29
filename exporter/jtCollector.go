package exporter

import (
	"fmt"
	"github.com/foreso-GitHub/jingtum-monitor/common"
	"github.com/prometheus/client_golang/prometheus"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
)

var config = common.LoadConfig("./config/config.json")

var nodeCount int32
var onlineNodeCount int32
var consensusNodeCount int32
var blockNumber int32

type JtCollector struct {
	nodeCountDesc          *prometheus.Desc
	onlineNodeCountDesc    *prometheus.Desc
	consensusNodeCountDesc *prometheus.Desc
	jtBlockNumberDesc      *prometheus.Desc
	guard                  sync.Mutex
}

func NewJtCollector() prometheus.Collector {
	return &JtCollector{
		nodeCountDesc: prometheus.NewDesc(
			"jt_total_node_count",
			"井通区块链网络中的节点数",
			nil, nil),
		onlineNodeCountDesc: prometheus.NewDesc(
			"jt_online_node_count",
			"井通区块链网络中在线的节点数",
			nil, nil),
		consensusNodeCountDesc: prometheus.NewDesc(
			"jt_consensus_node_count",
			"井通区块链网络中参与共识的节点数",
			nil, nil),
		jtBlockNumberDesc: prometheus.NewDesc(
			"jt_block_number",
			"井通区块链网络当前的区块高度",
			nil, nil),
	}
}

func (n *JtCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- n.nodeCountDesc
	ch <- n.onlineNodeCountDesc
	ch <- n.consensusNodeCountDesc
	ch <- n.jtBlockNumberDesc
}

func (n *JtCollector) Collect(ch chan<- prometheus.Metric) {
	network := LoadJtNetworkConfig(config.JtConfigPath)
	FlushNetwork(&network)
	n.guard.Lock()
	ch <- prometheus.MustNewConstMetric(n.nodeCountDesc, prometheus.GaugeValue, float64(network.NodeCount))
	ch <- prometheus.MustNewConstMetric(n.onlineNodeCountDesc, prometheus.GaugeValue, float64(network.OnlineNodeCount))
	ch <- prometheus.MustNewConstMetric(n.consensusNodeCountDesc, prometheus.GaugeValue, float64(network.ConsensusNodeCount))
	ch <- prometheus.MustNewConstMetric(n.jtBlockNumberDesc, prometheus.GaugeValue, float64(network.BlockNumber))
	n.guard.Unlock()
}

func GetBlockNumberInfo() []byte {
	url := "http://139.198.177.59:9545/v1/jsonrpc" //请求地址
	contentType := "application/json"
	data := strings.NewReader("{\"jsonrpc\":\"2.0\",\"method\":\"jt_blockNumber\",\"params\":[],\"id\":1}")
	resp, err := http.Post(url, contentType, data)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	//bodyString := string(body[:])
	//fmt.Printf("[GetBlockNumber] response:%s\n",bodyString)
	return body
}
