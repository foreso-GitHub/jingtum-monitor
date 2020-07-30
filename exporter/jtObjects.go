package exporter

import (
	"encoding/json"
	"fmt"
	"github.com/foreso-GitHub/jingtum-monitor/common"
	"io/ioutil"
	"log"
)

//region define jt objects

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
	Ip          string `json:"ip"`
	Port        string `json:"port"`
	Name        string `json:"name"`
	Online      bool
	Consensus   bool
	BlockNumber int
	LatestBlock JtBlock
}

type JtNetwork struct {
	Name               string `json:"name"`
	NodeCount          int
	NodeList           []JtNode `json:"nodes"`
	OnlineNodeCount    int
	OnlineRate         float32
	ConsensusNodeCount int
	ConsensusRate      float32
	BlockNumber        int
	LatestBlock        JtBlock
}

type JtResponseJson struct {
	Jsonrpc string `json:"jsonrpc"`
	Id      int    `json:"id"`
	Status  string `json:"status"`
}

type BlockNumberJson struct {
	JtResponseJson
	BlockNumber int `json:"result"`
}

type BlockJson struct {
	JtResponseJson
	Block JtBlock `json:"result"`
}

//endregion

func FlushNode(node *JtNode) {
	url := "http://" + node.Ip + ":" + node.Port + "/v1/jsonrpc" //请求地址

	blockNumber, err := GetBlockNumber(url)
	if err != nil {
		node.Online = false
	} else {
		node.Online = true
		node.BlockNumber = blockNumber
	}

	if block, err := GetBlockByNumber(url, blockNumber); err == nil {
		node.LatestBlock = *block
	}

	fmt.Println(node)
}

func FlushNetwork(network *JtNetwork) {
	network.NodeCount = len(network.NodeList)
	network.BlockNumber = -1
	network.OnlineNodeCount = 0
	network.ConsensusNodeCount = 0
	for i := 0; i < network.NodeCount; i++ {
		node := &network.NodeList[i]
		FlushNode(node)
		if network.BlockNumber < node.BlockNumber {
			network.BlockNumber = node.BlockNumber
			network.LatestBlock = node.LatestBlock
		}
		if node.Online {
			network.OnlineNodeCount++
		}
	}
	for i := 0; i < network.NodeCount; i++ {
		network.NodeList[i].Consensus = network.BlockNumber-network.NodeList[i].BlockNumber <= 2
		if network.NodeList[i].Consensus {
			network.ConsensusNodeCount++
		}
	}
	network.OnlineRate = float32(network.OnlineNodeCount) / float32(network.NodeCount) * 100
	network.ConsensusRate = float32(network.ConsensusNodeCount) / float32(network.NodeCount) * 100
	fmt.Println(network)
}

func LoadJtNetworkConfig(path string) JtNetwork {
	var network JtNetwork
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		log.Panicln("load config file failed: ", err)
	}
	err = json.Unmarshal(buf, &network)
	if err != nil {
		log.Panicln("decode config file failed:", string(buf), err)
	}
	return network
}

func GetUrlFromNode(node *JtNode) string {
	url := "http://" + node.Ip + ":" + node.Port + "/v1/jsonrpc" //请求地址
	return url
}

func PickNode(nodes []JtNode) *JtNode {
	count := len(nodes)
	if count > 0 {
		index := common.Rand(count)
		return &nodes[index]
	} else {
		return nil
	}
}

func GetRandUrl(nodes []JtNode) string {
	node := PickNode(nodes)
	url := GetUrlFromNode(node)
	_, error := GetBlockNumber(url)
	if error != nil {
		return GetRandUrl(nodes)
	}
	return url
}
