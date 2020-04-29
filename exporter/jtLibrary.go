package exporter

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

var m_ID = 1

//region jt library

func GetBlockNumber(url string, contentType string) ([]byte, error) {
	params := GenerateRequest("jt_blockNumber", "")
	return Post(url, contentType, params)
}

func GetBlockByNumber(url string, contentType string, blockNumber int) ([]byte, error) {

	params := GenerateRequest("jt_getBlockByNumber", "\""+strconv.Itoa(blockNumber)+"\",false")
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

	fmt.Println(request)
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
