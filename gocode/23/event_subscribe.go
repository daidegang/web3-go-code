package main

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
)

// 订阅事件日志
func main() {
	// 拨打启用websocket的以太坊客户端
	client, err := ethclient.Dial("ws://127.0.0.1:8545")
	if err != nil {
		log.Fatal(err)
	}

	// 加载合约地址
	contractAddress := common.HexToAddress("0xe7f1725e7734ce288f8367e1bb143e90bb3f0512")

	// 筛选查询合约事件
	query := ethereum.FilterQuery{
		Addresses: []common.Address{contractAddress},
	}

	// 接收事件
	logs := make(chan types.Log)
	// 客户端发起事件日志订阅
	sub, err := client.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		log.Fatal(err)
	}

	// 使用for-select语句连续循环读入新的日志
	for {
		select {
		case err := <-sub.Err():
			log.Fatal(err)
		case vLog := <-logs:
			fmt.Println(vLog)
		}
	}
}
