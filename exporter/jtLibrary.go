package exporter

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var m_ID = 1

const CONTENT_TYPE = "application/json"

//region jt library

func GetBlockNumber(url string) (int, error) {
	params := GenerateRequest("jt_blockNumber", "")
	jsonString, err := Post(url, CONTENT_TYPE, params)
	if err != nil {
		log.Println("[jt_blockNumber] error: ", err)
		return -1, err
	} else if string(jsonString[:]) == "Bad Request\n" {
		log.Println("[jt_blockNumber] error: "+url+" | ", "Bad Request")
		return -2, errors.New("Bad Request")
	} else {
		var blockNumberJson BlockNumberJson
		if err := json.Unmarshal(jsonString, &blockNumberJson); err == nil {
			//log.Println("================json str 转BlockNumberJson==")
			//log.Println("blockNumberJson: ", blockNumberJson)
			//log.Println(blockNumberJson.BlockNumber)
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
		log.Println("================[jt_getBlockByNumber] error==", url)
		log.Println(err)
		return nil, err
	} else if string(jsonString[:]) == "Bad Request\n" {
		log.Println("================[jt_getBlockByNumber] error==", url)
		log.Println("Bad Request")
		return nil, err
	}
	var blockJson BlockJson
	if err := json.Unmarshal(jsonString, &blockJson); err == nil {
		//log.Println("================json str 转BlockNumberJson==")
		//log.Println(blockNumberJson)
		//log.Println(blockNumberJson.BlockNumber)
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

	//fmt.Printf("===request: %+v\n", request)
	return request
}

//Post method with timeout
func Post(url string, contentType string, params string) ([]byte, error) {
	//start := time.Now()
	trans := &http.Transport{}
	client := &http.Client{
		Transport: trans,
		Timeout:   time.Duration(config.RequestTimeout) * time.Millisecond,
	}
	req, err := http.NewRequest("Post", url, bytes.NewBuffer([]byte(params)))
	if err != nil {
		panic(err)
	}

	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	//log.Println("CostTime is: " + time.Since(start).String())
	defer response.Body.Close()
	return ioutil.ReadAll(response.Body)
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
	//	//fmt.Printf("%v: abort\n", time.Now())
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
