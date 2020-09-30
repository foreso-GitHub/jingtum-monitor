package common

import (
	"encoding/json"
	"github.com/foreso-GitHub/jingtum-monitor/types"
	"io/ioutil"
	"log"
)

type ExporterConfig struct {
	JtConfigPath      string       `json:"jtConfigPath"`
	ExportAddress     string       `json:"exportAddress"`
	SupervisorMode    int          `json:"supervisorMode"` //1: Network mode, 2: LocalNode mode, 3: Both
	RequestTimeout    int          `json:"requestTimeout"` //must > 2000, otherwise prometheus will not fresh in time and then freeze.
	RequestRetrySpan  int          `json:"requestRetrySpan"`
	RequestRetryLimit int          `json:"requestRetryLimit"`
	LocalJtNode       types.JtNode `json:"localJtNode"`
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
