package exporter

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

var m_ID = 1

const CONTENT_TYPE = "application/json"

//region jt library

func GetBlockNumber(url string) (int, error) {
	params := GenerateRequest("jt_blockNumber", "")
	jsonString, err := Post(url, CONTENT_TYPE, params)
	if err != nil {
		fmt.Println("================[jt_blockNumber] error==", url)
		fmt.Println(err)
		return -1, err
	} else if string(jsonString[:]) == "Bad Request\n" {
		fmt.Println("================[jt_blockNumber] error==", url)
		fmt.Println("Bad Request")
		return -2, errors.New("Bad Request")
	} else {
		var blockNumberJson BlockNumberJson
		if err := json.Unmarshal(jsonString, &blockNumberJson); err == nil {
			//fmt.Println("================json str 转BlockNumberJson==")
			//fmt.Println(blockNumberJson)
			//fmt.Println(blockNumberJson.BlockNumber)
			return blockNumberJson.BlockNumber, nil
		} else {
			return -3, errors.New("Unmarshal error")
		}
	}
}

func GetBlockByNumber(url string, blockNumber int) (*JtBlock, error) {
	params := GenerateRequest("jt_getBlockByNumber", "\""+strconv.Itoa(blockNumber)+"\",false")
	jsonString, err := Post(url, CONTENT_TYPE, params)

	if err != nil {
		fmt.Println("================[jt_getBlockByNumber] error==", url)
		fmt.Println(err)
		return nil, err
	} else if string(jsonString[:]) == "Bad Request\n" {
		fmt.Println("================[jt_getBlockByNumber] error==", url)
		fmt.Println("Bad Request")
		return nil, err
	}
	var blockJson BlockJson
	if err := json.Unmarshal(jsonString, &blockJson); err == nil {
		//fmt.Println("================json str 转BlockNumberJson==")
		//fmt.Println(blockNumberJson)
		//fmt.Println(blockNumberJson.BlockNumber)
		return &blockJson.Block, nil
	} else {
		return nil, errors.New("Unmarshal error")
	}
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

	fmt.Printf("===request: %+v\n", request)
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
