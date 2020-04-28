package exporter

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

var m_ID = 1

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

//endregion

func FlushNode(node *JtNode) {
	url := "http://" + node.Ip + ":" + node.Port + "/v1/jsonrpc" //请求地址
	contentType := "application/json"
	jsonString, err := GetBlockNumber(url, contentType)
	if err != nil {
		fmt.Println("================[jt_blockNumber] error==")
		fmt.Println(err)
		node.Online = false
	} else {
		node.Online = true
	}
	var blockNumberJson BlockNumberJson
	if err := json.Unmarshal(jsonString, &blockNumberJson); err == nil {
		//fmt.Println("================json str 转BlockNumberJson==")
		//fmt.Println(blockNumberJson)
		//fmt.Println(blockNumberJson.BlockNumber)
		node.BlockNumber = blockNumberJson.BlockNumber
	} else {
		node.BlockNumber = -1
	}

	fmt.Println(node)

	//params_getBlockNumber := strings.NewReader("{\"jsonrpc\":\"2.0\",\"method\":\"jt_blockNumber\",\"params\":[],\"id\":1}")

}

func FlushNetwork() JtNetwork {
	var network JtNetwork
	network.NodeCount = 3
	network.BlockNumber = -1
	network.OnlineNodeCount = 0
	network.ConsensusNodeCount = 0
	network.NodeList = make([]JtNode, network.NodeCount)

	for i := 0; i < network.NodeCount; i++ {
		var node JtNode
		node.Name = "Node_" + strconv.Itoa(i+1)
		node.Ip = "139.198.177.59"
		node.Port = "9545"
		FlushNode(&node)
		network.NodeList[i] = node

		if network.BlockNumber < node.BlockNumber {
			network.BlockNumber = node.BlockNumber
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

	return network
}

//region jt library

func GetBlockNumber(url string, contentType string) ([]byte, error) {
	params := GenerateRequest("jt_blockNumber", "")
	return Post(url, contentType, params)
}

func GetBlockByNumber(url string, contentType string, blockNumber int) ([]byte, error) {
	//params_getBlockNumber := "{\"jsonrpc\":\"2.0\",\"method\":\"jt_getBlockByNumber\"," +
	//	"\"params\":[" + strconv.Itoa(blockNumber) + ",false],\"id\":1}"
	params := GenerateRequest("jt_getBlockByNumber", strconv.Itoa(blockNumber)+",false")
	return Post(url, contentType, params)
}

func GenerateRequest(method string, params string) string {
	request := "{\"jsonrpc\":\"2.0\",\"id\":"
	request += strconv.Itoa(m_ID)
	m_ID += 1

	request += ",\"method\":\""
	request += method

	request += "\",\"params\":["
	request += params

	request += "]}"

	//fmt.Println(request)
	return request
}

func Post(url string, contentType string, params string) ([]byte, error) {
	reader := strings.NewReader(params)
	resp, err := http.Post(url, contentType, reader)
	if err != nil {
		//fmt.Println("================[jt_blockNumber] error==")
		//fmt.Println(err)
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

//endregion
