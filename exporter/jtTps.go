package exporter

import (
	"errors"
	"github.com/foreso-GitHub/jingtum-monitor/types"
	"log"
	"strconv"
)

const BLOCK_PERIOD = 5 // system generate 1 block per 5 seconds

//region init

func CreateJtTpsStatus(initBlockNumber int) *types.JtTpsStatus {
	status := new(types.JtTpsStatus)
	status.CurrentBlockNumber = initBlockNumber
	status.BlockMap = make(map[int]types.JtBlock)
	status.Blocks = make([]types.JtBlock, 0)
	status.TotalBlockCount = 0
	status.TotalPeriod = 0
	status.TotalTxCount = 0
	status.TotalTps = 0
	status.TpsMap = make(map[int]types.JtTps)
	AddJtTps("最新单个区块TPS", 1, status)
	AddJtTps("最新三个区块TPS", 3, status)
	AddJtTps("最近一分钟TPS", 1*12, status)
	AddJtTps("最近一小时TPS", 1*12*60, status)
	AddJtTps("最近一天TPS", 1*12*60*24, status)
	AddJtTps("最近一周TPS", 1*12*60*24*7, status)
	//log.Println("status: %+v\n", status)
	return status
}

func InitJtTps(name string, blockCount int) *types.JtTps {
	tps := new(types.JtTps)
	tps.Name = name
	tps.BlockCount = blockCount
	tps.Period = tps.BlockCount * BLOCK_PERIOD
	tps.TxCount = 0
	tps.Tps = 0
	tps.Blocks = make([]types.JtBlock, 0)
	//log.Println("tps: %+v\n", tps)
	return tps
}

func AddJtTps(name string, blockCount int, status *types.JtTpsStatus) {
	tps := InitJtTps(name, blockCount)
	status.TpsMap[blockCount] = *tps
}

//endregion

//region flush

func FlushTpsStatus(status *types.JtTpsStatus) bool {
	_, newblockNumber, err := GetBlockNumberByRandNode()
	if err != nil {
		return false
	} else {
		lastBlockNumber := status.CurrentBlockNumber
		log.Println("BlockNumber: ", lastBlockNumber, " - ", newblockNumber)

		if newblockNumber > lastBlockNumber {
			for blockNumber := lastBlockNumber + 1; blockNumber <= newblockNumber; blockNumber++ {
				if _, block, err := GetBlockByNumberByRandNode(blockNumber); err == nil {
					txCount := len(block.Transactions)
					log.Println("blockNumber: ", blockNumber, " | tx count: ", txCount)

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
						for key, _ := range status.TpsMap {
							tps := status.TpsMap[key]
							blockCount := tps.BlockCount
							if blockCount <= status.TotalBlockCount {
								start := status.TotalBlockCount - blockCount
								end := status.TotalBlockCount
								flushBlocks := status.Blocks[start:end]
								FlushSingleTps(&tps, flushBlocks)
								status.TpsMap[key] = tps
							}
						}
					}
				} else {
					log.Println("GetBlockByNumber error: %+v\n", err)
					return false
				}
			}
		} else {
			return true
		}
	}
	return true
}

func FlushSingleTps(tps *types.JtTps, blocks []types.JtBlock) error {
	if tps.BlockCount != len(blocks) {
		return errors.New("The count of blocks doesn't match!  The correct count should be " + strconv.Itoa(tps.BlockCount))
	}
	//tps.Blocks = blocks  //todo need restore later, now remove just for debug
	txCount := 0
	for i := 0; i < tps.BlockCount; i++ {
		txCount += len(blocks[i].Transactions)
	}
	tps.TxCount = txCount
	tps.Tps = float64(tps.TxCount) / float64(tps.Period)
	log.Println("===flush tps: ", tps)
	return nil
}

//endregion
