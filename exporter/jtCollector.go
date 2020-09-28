package exporter

import (
	"github.com/foreso-GitHub/jingtum-monitor/common"
	"github.com/prometheus/client_golang/prometheus"
	"sync"
)

var config = common.LoadConfig("./config/config.json")
var firstCollect = true
var tpsStatus = CreateJtTpsStatus(-1)

var nodeCount int32
var onlineNodeCount int32
var consensusNodeCount int32
var blockNumber int32

var blockTps float32
var minuteTps float32
var hourTps float32
var dayTps float32
var weekTps float32
var highestTps float32
var totalTps float32

//var blockTxCount int32

type JtCollector struct {
	nodeCountDesc          *prometheus.Desc
	onlineNodeCountDesc    *prometheus.Desc
	consensusNodeCountDesc *prometheus.Desc
	jtBlockNumberDesc      *prometheus.Desc
	jtBlockTxCountDesc     *prometheus.Desc
	blockTpsDesc           *prometheus.Desc
	minuteTpsDesc          *prometheus.Desc
	hourTpsDesc            *prometheus.Desc
	dayTpsDesc             *prometheus.Desc
	weekTpsDesc            *prometheus.Desc
	totalTpsDesc           *prometheus.Desc
	highestTpsDesc         *prometheus.Desc
	guard                  sync.Mutex
}

func NewJtCollector() prometheus.Collector {
	return &JtCollector{
		nodeCountDesc: prometheus.NewDesc(
			"jt_total_node_count",
			"井通区块链网络中的节点数",
			nil, nil),
		onlineNodeCountDesc: prometheus.NewDesc(
			"jt_online_node_count",
			"井通区块链网络中在线的节点数",
			nil, nil),
		consensusNodeCountDesc: prometheus.NewDesc(
			"jt_consensus_node_count",
			"井通区块链网络中参与共识的节点数",
			nil, nil),
		jtBlockNumberDesc: prometheus.NewDesc(
			"jt_block_number",
			"井通区块链网络当前的区块高度",
			nil, nil),
		jtBlockTxCountDesc: prometheus.NewDesc(
			"jt_current_block_tx_count",
			"井通区块链网络当前区块交易数",
			nil, nil),
		blockTpsDesc: prometheus.NewDesc(
			"jt_current_block_tps",
			"井通区块链网络最新区块TPS",
			nil, nil),
		minuteTpsDesc: prometheus.NewDesc(
			"jt_current_minute_tps",
			"井通区块链网络最近一分钟TPS",
			nil, nil),
		hourTpsDesc: prometheus.NewDesc(
			"jt_current_hour_tps",
			"井通区块链网络最近一小时TPS",
			nil, nil),
		dayTpsDesc: prometheus.NewDesc(
			"jt_current_day_tps",
			"井通区块链网络最近一天TPS",
			nil, nil),
		weekTpsDesc: prometheus.NewDesc(
			"jt_current_week_tps",
			"井通区块链网络最近一周TPS",
			nil, nil),
		totalTpsDesc: prometheus.NewDesc(
			"jt_average_tps",
			"井通区块链网络平均TPS",
			nil, nil),
		highestTpsDesc: prometheus.NewDesc(
			"jt_highest_tps",
			"井通区块链网络峰值TPS",
			nil, nil),
	}
}

func (n *JtCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- n.nodeCountDesc
	ch <- n.onlineNodeCountDesc
	ch <- n.consensusNodeCountDesc
	ch <- n.jtBlockNumberDesc
	ch <- n.jtBlockTxCountDesc

	ch <- n.blockTpsDesc
	ch <- n.minuteTpsDesc
	ch <- n.hourTpsDesc
	ch <- n.dayTpsDesc
	ch <- n.weekTpsDesc
	ch <- n.totalTpsDesc
	ch <- n.highestTpsDesc
}

func (n *JtCollector) Collect(ch chan<- prometheus.Metric) {
	//flush blockchain network
	network := LoadJtNetworkConfig(config.JtConfigPath)
	FlushNetwork(&network)

	//flush tps
	if firstCollect {
		_, blockNumber, _ := GetBlockNumberByRandNode()
		tpsStatus = CreateJtTpsStatus(blockNumber)
		firstCollect = false
	}
	_ = FlushTpsStatus(tpsStatus)
	//flushOK := FlushTpsStatus(tpsStatus)
	//log.Println("flushOK: %+v\n", flushOK)

	n.guard.Lock()
	ch <- prometheus.MustNewConstMetric(n.nodeCountDesc, prometheus.GaugeValue, float64(network.NodeCount))
	ch <- prometheus.MustNewConstMetric(n.onlineNodeCountDesc, prometheus.GaugeValue, float64(network.OnlineNodeCount))
	ch <- prometheus.MustNewConstMetric(n.consensusNodeCountDesc, prometheus.GaugeValue, float64(network.ConsensusNodeCount))
	ch <- prometheus.MustNewConstMetric(n.jtBlockNumberDesc, prometheus.GaugeValue, float64(network.BlockNumber))
	ch <- prometheus.MustNewConstMetric(n.jtBlockTxCountDesc, prometheus.GaugeValue, float64(len(network.LatestBlock.Transactions)))

	ch <- prometheus.MustNewConstMetric(n.blockTpsDesc, prometheus.GaugeValue, tpsStatus.TpsMap[1].Tps)
	ch <- prometheus.MustNewConstMetric(n.minuteTpsDesc, prometheus.GaugeValue, tpsStatus.TpsMap[1*12].Tps)
	ch <- prometheus.MustNewConstMetric(n.hourTpsDesc, prometheus.GaugeValue, tpsStatus.TpsMap[1*12*60].Tps)
	ch <- prometheus.MustNewConstMetric(n.dayTpsDesc, prometheus.GaugeValue, tpsStatus.TpsMap[1*12*60*24].Tps)
	ch <- prometheus.MustNewConstMetric(n.weekTpsDesc, prometheus.GaugeValue, tpsStatus.TpsMap[1*12*60*24*7].Tps)
	ch <- prometheus.MustNewConstMetric(n.totalTpsDesc, prometheus.GaugeValue, tpsStatus.TotalTps)
	n.guard.Unlock()
}
