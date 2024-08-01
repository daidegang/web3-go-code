package main

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math/big"
)

// 查询区块信息
func main() {
	client, err := ethclient.Dial("https://cloudflare-eth.com")
	if err != nil {
		log.Fatal(err)
	}

	header, err := client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}
	//区块头信息
	fmt.Printf("header number: %v\n", header.Number.String())

	blockNumber := big.NewInt(5671744)
	block, err := client.BlockByNumber(context.Background(), blockNumber)
	if err != nil {
		log.Fatal(err)
	}
	//区块号
	fmt.Printf("block number: %v\n", block.Number().Uint64())
	//区块时间戳
	fmt.Printf("block time: %v\n", block.Time())
	//区块摘要
	fmt.Printf("block difficulty: %v\n", block.Difficulty().Uint64())
	//区块难度
	fmt.Printf("block hash: %v\n", block.Hash().Hex())
	//交易列表
	fmt.Printf("block transactions: %v\n", block.Transactions())

	//获取交易数目
	count, err := client.TransactionCount(context.Background(), block.Hash())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("number of transactions: %d\n", count)
}
