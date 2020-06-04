package exporter

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
)

const BLOCK_PERIOD = 5

//region tps struct

type JtTps struct {
	Name       string
	Period     int
	BlockCount int
	Blocks     []JtBlock
	TxCount    int
	Tps        float64
}

type JtTpsStatus struct {
	CurrentBlockNumber int
	BlockMap           map[int]JtBlock
	Blocks             []JtBlock
	TotalBlockCount    int
	TotalPeriod        int
	TotalTxCount       int
	TotalTps           float64
	TpsMap             map[int]JtTps //key = 1, 12, 12*60, 12*60*24, 12*60*24*7, total
}

//endregion

//region init

func InitJtTpsStatus() *JtTpsStatus {
	status := new(JtTpsStatus)
	status.CurrentBlockNumber = -1
	status.BlockMap = make(map[int]JtBlock)
	status.Blocks = make([]JtBlock, 0)
	status.TotalBlockCount = 0
	status.TotalPeriod = 0
	status.TotalTxCount = 0
	status.TotalTps = 0
	status.TpsMap = make(map[int]JtTps)
	AddJtTps("最新区块TPS", 1, status)
	AddJtTps("最近一分钟TPS", 1*12, status)
	AddJtTps("最近一小时TPS", 1*12*60, status)
	AddJtTps("最近一天TPS", 1*12*60*24, status)
	AddJtTps("最近一周TPS", 1*12*60*24*7, status)
	fmt.Println("status: %+v", status)
	return status
}

func InitJtTps(name string, blockCount int) *JtTps {
	tps := new(JtTps)
	tps.Name = name
	tps.BlockCount = blockCount
	tps.Period = tps.BlockCount * BLOCK_PERIOD
	tps.TxCount = 0
	tps.Tps = 0
	tps.Blocks = make([]JtBlock, 0)
	fmt.Println("tps: %+v", tps)
	return tps
}

func AddJtTps(name string, blockCount int, status *JtTpsStatus) {
	tps := InitJtTps(name, blockCount)
	status.TpsMap[blockCount] = *tps
}

//endregion

//region flush

func FlushTpsStatus(url string, status *JtTpsStatus) bool {
	newblockNumber, err := GetBlockNumber(url)
	if err != nil {
		return false
	} else {
		lastBlockNumber := status.CurrentBlockNumber

		//todo: for test only
		gap := 3
		if newblockNumber > gap {
			lastBlockNumber = newblockNumber - gap
		}

		fmt.Println("lastBlockNumber: %+v", lastBlockNumber)
		fmt.Println("newblockNumber: %+v", newblockNumber)

		if newblockNumber > lastBlockNumber {
			for blockNumber := lastBlockNumber + 1; blockNumber <= newblockNumber; blockNumber++ {
				if block, err := GetBlockByNumber(url, blockNumber); err == nil {
					//txCount := len(block.Transactions)
					txCount := rand.Intn(100)         //todo: fake tx count, need be deleted later.
					block.Parent_close_time = txCount //todo: use Parent_close_time to transfer fake tx count, need be deleted later.
					fmt.Println("blockNumber: %+v", blockNumber)
					fmt.Println("tx count: %+v", txCount)

					block := *block
					status.Blocks = append(status.Blocks, block)
					status.BlockMap[blockNumber] = block

					status.TotalBlockCount++
					status.TotalPeriod += BLOCK_PERIOD
					status.TotalTxCount += txCount
					status.TotalTps = float64(status.TotalTxCount) / float64(status.TotalPeriod)
					status.CurrentBlockNumber = blockNumber

					//reflush tps list, only flush on last block
					if blockNumber == newblockNumber {
						for _, tps := range status.TpsMap {
							blockCount := tps.BlockCount
							if blockCount <= status.TotalBlockCount {
								start := status.TotalBlockCount - blockCount
								end := status.TotalBlockCount
								flushBlocks := status.Blocks[start:end]
								FlushSingleTps(tps, flushBlocks)
							}
						}
					}
				} else {
					return false
				}
			}
		} else {
			return true
		}
	}
	return true
}

func FlushSingleTps(tps JtTps, blocks []JtBlock) error {
	if tps.BlockCount != len(blocks) {
		return errors.New("The count of blocks doesn't match!  The correct count should be " + strconv.Itoa(tps.BlockCount))
	}
	tps.Blocks = make([]JtBlock, 0)
	for i := 0; i < tps.BlockCount; i++ {
		//tps.Blocks = append(tps.Blocks, blocks[i])
		tps.TxCount += len(blocks[i].Transactions)
		tps.TxCount += blocks[i].Parent_close_time
	}
	tps.Tps = float64(tps.TxCount) / float64(tps.Period)
	fmt.Println("===flush tps: %+v", tps)
	return nil
}

//endregion
