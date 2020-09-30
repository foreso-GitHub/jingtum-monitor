package exporter

import (
	"encoding/json"
	"github.com/foreso-GitHub/jingtum-monitor/common"
	"github.com/foreso-GitHub/jingtum-monitor/types"
	"io/ioutil"
	"log"
)

func FlushNode(node *types.JtNode) {
	url := "http://" + node.Ip + ":" + node.Port + "/v1/jsonrpc" //请求地址

	blockNumber, err := GetBlockNumberByNode(url)
	if err != nil {
		node.Online = false
	} else {
		node.Online = true
		node.BlockNumber = blockNumber
	}
	//log.Println("===node: %+v\n", node)
}

func FlushNetwork(network *types.JtNetwork) {
	network.NodeCount = len(network.NodeList)
	network.BlockNumber = -1
	network.OnlineNodeCount = 0
	network.ConsensusNodeCount = 0
	for i := 0; i < network.NodeCount; i++ {
		node := &network.NodeList[i]
		FlushNode(node)
		if network.BlockNumber < node.BlockNumber {
			network.BlockNumber = node.BlockNumber
			_, block, _ := GetBlockByNumberByRandNode(network.BlockNumber)
			network.LatestBlock = *block
		}
		if node.Online {
			network.OnlineNodeCount++
		}
	}
	for i := 0; i < network.NodeCount; i++ {
		network.NodeList[i].Consensus = network.NodeList[i].Online && network.BlockNumber-network.NodeList[i].BlockNumber <= 2
		if network.NodeList[i].Consensus {
			network.ConsensusNodeCount++
		}
	}
	network.OnlineRate = float32(network.OnlineNodeCount) / float32(network.NodeCount) * 100
	network.ConsensusRate = float32(network.ConsensusNodeCount) / float32(network.NodeCount) * 100
	log.Println("===network: ", network)
}

func LoadJtNetworkConfig(path string) types.JtNetwork {
	var network types.JtNetwork
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

func GetUrlFromNode(node *types.JtNode) string {
	url := "http://" + node.Ip + ":" + node.Port + "/v1/jsonrpc" //请求地址
	return url
}

func PickNode(nodes []types.JtNode) *types.JtNode {
	count := len(nodes)
	if count > 0 {
		index := common.Rand(count)
		return &nodes[index]
	} else {
		return nil
	}
}

func GetRandUrl(nodes []types.JtNode) string {
	node := PickNode(nodes)
	url := GetUrlFromNode(node)
	//_, error := GetBlockNumber(url)
	//if error != nil {
	//	time.Sleep(100 * time.Millisecond)
	//	return GetRandUrl(nodes)
	//}
	return url
}
