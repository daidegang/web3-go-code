package main

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math/big"
	"strings"
	store "web3-go-code/contract/bindings"
)

// 读取事件日志
func main() {
	// 拨打启用websocket的以太坊客户端
	client, err := ethclient.Dial("ws://127.0.0.1:8545")
	if err != nil {
		log.Fatal(err)
	}

	// 加载合约地址
	contractAddress := common.HexToAddress("0xe7f1725e7734ce288f8367e1bb143e90bb3f0512")
	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(2394201),
		ToBlock:   big.NewInt(2394201),
		Addresses: []common.Address{contractAddress},
	}

	// 接收查询并将返回的匹配事件日志
	logs, err := client.FilterLogs(context.Background(), query)
	if err != nil {
		log.Fatal(err)
	}

	// 导入智能合约ABI编码
	contractAbi, err := abi.JSON(strings.NewReader(string(store.StoreMetaData.ABI)))
	if err != nil {
		log.Fatal(err)
	}

	// 通过日志进行迭代并进行解码
	for _, vLog := range logs {
		fmt.Println(vLog.BlockHash.Hex())
		fmt.Println(vLog.BlockNumber)
		fmt.Println(vLog.TxHash.Hex())

		event := struct {
			Key   [32]byte
			Value [32]byte
		}{}

		err := contractAbi.UnpackIntoInterface(&event, "ItemSet", vLog.Data)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(event.Key[:]))
		fmt.Println(string(event.Value[:]))

		// 读取事件主题
		var topics [4]string
		for i := range vLog.Topics {
			topics[i] = vLog.Topics[i].Hex()
		}
		fmt.Println(topics[0])
	}

	eventSignature := []byte("ItemSet(byte32,byte32)")
	hash := crypto.Keccak256Hash(eventSignature)
	fmt.Println(hash.Hex())
}
