package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rlp"
	"log"
)

// 发送原始交易事务
func main() {
	client, err := ethclient.Dial("http://127.0.0.1:8545")
	if err != nil {
		log.Fatal(err)
	}

	// 首先将原始事务十六进制解码为字节格式
	rawTx := "f86ff86d0184693fa2cd8252089470997970c51812dc3a010c7d01b50e0d17dc79c888016345785d8a00008082f4f6a0630f7d055601b1b2353f5146ef43cd70bc6bcaee10e4333ca7039296db15c3e8a058b80f79e8b0c91a3956a0e7088241c9570537f34d3de549d24aa29fb599a052"
	rawTxBytes, err := hex.DecodeString(rawTx)

	// 将原始事务字节和指针传递给以太坊事务类型
	tx := new(types.Transaction)
	rlp.DecodeBytes(rawTxBytes, tx)

	// 广播交易
	err = client.SendTransaction(context.Background(), tx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Transaction sent: %s", tx.Hash().Hex())
}
