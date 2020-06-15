package collector

import (
	"encoding/json"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
)

var nodeCount int32
var consensusNodeCount int32
var blockNumber int32

type JtCollector struct {
	nodeCountDesc          *prometheus.Desc
	consensusNodeCountDesc *prometheus.Desc
	jtBlockNumberDesc      *prometheus.Desc
	guard                  sync.Mutex
}

type JtBlock struct {
	Accepted              bool     `json:"accepted"`
	Account_hash          string   `json:"account_hash"`
	Close_flags           int      `json:"close_flags"`
	Close_time            int      `json:"close_time"`
	Close_time_human      string   `json:"close_time_human"`
	Close_time_resolution int      `json:"close_time_resolution"`
	Closed                bool     `json:"closed"`
	Hash                  string   `json:"hash"`
	Ledger_hash           string   `json:"ledger_hash"`
	Ledger_index          string   `json:"ledger_index"`
	Parent_close_time     int      `json:"parent_close_time"`
	Parent_hash           string   `json:"parent_hash"`
	SeqNum                string   `json:"seqNum"`
	TotalCoins            string   `json:"totalCoins"`
	Total_coins           string   `json:"total_coins"`
	Transaction_hash      string   `json:"transaction_hash"`
	Transactions          []string `json:"transactions"`
}

type JtNode struct {
	Ip          string
	Port        string
	Name        string
	Online      bool
	Consensus   bool
	BlockNumber int
	LatestBlock JtBlock
}

type JtNetwork struct {
	NodeCount          int
	NodeList           []JtNode
	OnlineNodeCount    int
	OnlineRate         float32
	ConsensusNodeCount int
	ConsensusRate      float32
	BlockNumber        int
	LatestBlock        JtBlock
}

type BlockNumberJson struct {
	Jsonrpc     string `json:"jsonrpc"`
	Id          int    `json:"id"`
	Status      string `json:"status"`
	BlockNumber int    `json:"result"`
}

func NewJtCollector() prometheus.Collector {
	return &JtCollector{
		nodeCountDesc: prometheus.NewDesc(
			"jt_total_node_count",
			"井通区块链网络中的节点数",
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
	ch <- n.consensusNodeCountDesc
	ch <- n.jtBlockNumberDesc
}

func (n *JtCollector) Collect(ch chan<- prometheus.Metric) {
	blockNumberJsonString := GetBlockNumber()
	fmt.Printf("[GetBlockNumber] blockNumberJsonString:%s\n", blockNumberJsonString)
	var blockNumberJson BlockNumberJson
	if err := json.Unmarshal(blockNumberJsonString, &blockNumberJson); err == nil {
		fmt.Println("================json str 转struct==")
		fmt.Println(blockNumberJson)
		fmt.Println(blockNumberJson.BlockNumber)
	}

	n.guard.Lock()
	ch <- prometheus.MustNewConstMetric(n.nodeCountDesc, prometheus.GaugeValue, 5)
	consensusNodeCount := 0
	if blockNumberJson.BlockNumber >= 0 {
		consensusNodeCount = 1
	}
	ch <- prometheus.MustNewConstMetric(n.consensusNodeCountDesc, prometheus.GaugeValue, float64(consensusNodeCount))
	ch <- prometheus.MustNewConstMetric(n.jtBlockNumberDesc, prometheus.GaugeValue, float64(blockNumberJson.BlockNumber))
	n.guard.Unlock()
}

func GetBlockNumber() []byte {
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
