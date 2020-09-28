package common

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type ExporterConfig struct {
	JtConfigPath      string `json:"jtConfigPath"`
	ExportAddress     string `json:"exportAddress"`
	SupervisorMode    int    `json:"supervisorMode"`
	RequestTimeout    int    `json:"requestTimeout"`
	RequestRetrySpan  int    `json:"requestRetrySpan"`
	RequestRetryLimit int    `json:"requestRetryLimit"`
}

func LoadConfig(path string) ExporterConfig {
	var config ExporterConfig
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		log.Panicln("load config file failed: ", err)
	}
	err = json.Unmarshal(buf, &config)
	if err != nil {
		log.Panicln("decode config file failed:", string(buf), err)
	}
	return config
}
