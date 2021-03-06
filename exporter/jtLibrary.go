package exporter

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/foreso-GitHub/jingtum-monitor/types"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

//region jt library
const CONTENT_TYPE = "application/json"

//region get block number

func GetBlockNumberByNode(url string) (int, error) {
	blockNumber, err := GetBlockNumber(url)
	return blockNumber.(int), err
}

func GetBlockNumber(url string, args ...interface{}) (interface{}, error) {
	var param types.JtGetBlockNumberParam
	param.Type = "number"
	paramBytes, err := json.Marshal(param)
	params := GenerateRequest("jt_blockNumber", string(paramBytes))

	jsonString, err := Post(url, CONTENT_TYPE, params)
	if err != nil {
		log.Println("[jt_blockNumber] error: ", err)
		return GetJtErrorCode(err.Error()), err
	}

	var blockNumberResponse types.JtGetBlockNumberResponse
	if err := json.Unmarshal(jsonString, &blockNumberResponse); err == nil {
		blockNumber := blockNumberResponse.Results[0].BlockNumber
		return blockNumber, nil
	} else {
		log.Println("[jt_blockNumber] error: ", err, " | ", string(jsonString))
		return -104, err
	}
}

func GetBlockNumberByRandNode() (string, int, error) {
	url, blockNumber, err := getJtInfo(GetBlockNumber)
	return url, blockNumber.(int), err
}

//endregion

//region get block by number

func GetBlockByNumber(url string, args ...interface{}) (interface{}, error) {
	var block types.JtBlock
	blockNumber := (args[0].([]interface{}))[0].(int)
	if blockNumber < 0 {
		err := errors.New("BlockNumber must be larger than 0!")
		log.Println("[jt_getBlockByNumber] error: ", err)
		return &block, err
	}

	var param types.JtGetBlockParam
	param.Number = strconv.Itoa(blockNumber)
	param.Full = false
	paramBytes, err := json.Marshal(param)
	params := GenerateRequest("jt_getBlockByNumber", string(paramBytes))

	jsonBytes, err := Post(url, CONTENT_TYPE, params)
	if err != nil {
		log.Println("[jt_getBlockByNumber] error: ", err)
		return &block, err
	}

	var blockResponse types.JtGetBlockResponse
	if err := json.Unmarshal(jsonBytes, &blockResponse); err == nil {
		blockResult := blockResponse.Results[0]
		block := blockResult.Block
		return &block, nil
	} else {
		log.Println("[jt_getBlockByNumber] error: ", err, " | ", string(jsonBytes))
		return &block, err
	}
}

func GetBlockByNumberByRandNode(blockNumber int) (string, *types.JtBlock, error) {
	url, block, err := getJtInfo(GetBlockByNumber, blockNumber)
	return url, block.(*types.JtBlock), err
}

//endregion

//endregion

//region common

func getJtInfo(jtFunction interface{}, jtFunctionArgs ...interface{}) (string, interface{}, error) {
	network := LoadJtNetworkConfig(config.JtConfigPath)
	nodes := network.NodeList
	return getJtInfoByNodes(nodes, config.RequestRetryLimit, 0, jtFunction, jtFunctionArgs)
}

var connectedUrl = ""

func getJtInfoByNodes(nodes []types.JtNode, retryLimit int, retriedCount int,
	jtFunction interface{}, jtFunctionArgs ...interface{}) (string, interface{}, error) {
	var url = ""
	if connectedUrl == "" {
		url = GetRandUrl(nodes)
		connectedUrl = url
	} else {
		url = connectedUrl
	}
	//log.Println("Url: %v\n", url)
	info, err := jtFunction.(func(string, ...interface{}) (interface{}, error))(url, jtFunctionArgs...)
	if err == nil || retriedCount == retryLimit || GetJtErrorCode(err.Error()) != -102 {
		return url, info, err
	} else {
		time.Sleep(time.Duration(config.RequestRetrySpan) * time.Millisecond)
		retriedCount++
		log.Println("Request retry count: ", retriedCount)
		connectedUrl = "" //reset connected url
		return getJtInfoByNodes(nodes, retryLimit, retriedCount, jtFunction, jtFunctionArgs...)
	}
}

func GetJtErrorCode(errString string) int {
	if errString != "" {
		if strings.Index(errString, "net/http: request canceled while waiting for connection") != -1 {
			return -102 //request timeout
		} else if errString == "Bad Request\n" {
			return -103 //bad request
		} else {
			return -101 //common request error
		}
	}
	return 0
}

//region request
var m_ID = 1

func GenerateRequest(method string, params string) string {
	request := "{\"jsonrpc\":\"2.0\",\"id\":"
	request += strconv.Itoa(m_ID)
	m_ID += 1

	request += ",\"method\":\""
	request += method

	request += "\",\"params\":["
	request += params

	request += "]}"

	//log.Println("===request: %+v\n", request)
	return request
}

//Post method with timeout
func Post(url string, contentType string, params string) ([]byte, error) {
	//start := time.Now()
	//log.Println("start is: ", start)
	trans := &http.Transport{}
	client := &http.Client{
		Transport: trans,
		Timeout:   time.Duration(config.RequestTimeout) * time.Millisecond,
	}

	request, err := http.NewRequest("Post", url, bytes.NewBuffer([]byte(params)))
	if err != nil {
		return nil, err
	}

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	//log.Println("CostTime is: " + time.Since(start).String())
	defer response.Body.Close()
	bytes, err := ioutil.ReadAll(response.Body)

	if string(bytes[:]) == "Bad Request\n" {
		return nil, errors.New("Bad Request")
	} else {
		return bytes, err
	}
}

//region deprecated post methods

//the first version of Post method
func Post_1(url string, contentType string, params string) ([]byte, error) {
	reader := strings.NewReader(params)
	resp, err := http.Post(url, contentType, reader)
	if err != nil {
		//log.Println("================[jt_blockNumber] error==")
		//log.Println(err)
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

//bad Post method which is not implemented well.
func Post2(url string, contentType string, params string) ([]byte, error) {
	start := time.Now()

	timeout := 10 * time.Millisecond
	context, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	//context, cancel := context.WithCancel(context.Background())
	//context, cancel := context.WithCancel(context.TODO())

	request, _ := http.NewRequest("POST", url, bytes.NewBuffer([]byte(params)))
	request.WithContext(context)

	//go func() {
	//	//time.Sleep(time.Second * time.Duration(1))
	//	//log.Println("%v: abort\n", time.Now())
	//	trans.CancelRequest(req)
	//}()

	log.Println("Request: Post | " + url + " | " + params)
	client := &http.Client{}
	resp, err := client.Do(request)

	if err != nil {
		fmt.Errorf("client.Do Error: %s", err.Error())
		return nil, err
	}

	log.Println("resp is: ", resp)
	log.Println("CostTime is: " + time.Since(start).String())
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

//endregion

//endregion

//endregion
