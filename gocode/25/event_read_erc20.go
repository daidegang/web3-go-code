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
	token "web3-go-code/contract/bindings"
)

// LogTransfer
type LogTransfer struct {
	From   common.Address
	To     common.Address
	Tokens *big.Int
}

// LogApproval
type LogApproval struct {
	TokenOwner common.Address
	Spender    common.Address
	Tokens     *big.Int
}

// 读取ERC-20代币事件日志
func main() {
	// 链接客户端
	client, err := ethclient.Dial("http://127.0.0.1:8545")
	if err != nil {
		log.Fatal(err)
	}

	// 加载合约地址
	contractAddress := common.HexToAddress("0x9fe46736679d2d9a65f0992f2272de9f3c7fa6e0")

	// 构建查询指定区块
	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(0),
		ToBlock:   big.NewInt(20),
		Addresses: []common.Address{
			contractAddress,
		},
	}

	// 查询合约事件日志
	logs, err := client.FilterLogs(context.Background(), query)
	if err != nil {
		log.Fatal(err)
	}

	// 加载合约ABI文件
	contractAbi, err := abi.JSON(strings.NewReader(string(token.TokenMetaData.ABI)))
	if err != nil {
		log.Fatal(err)
	}

	// 获取事情日志签名
	logTransferSig := []byte("Transfer(address,address,uint)")
	logApprovalSig := []byte("Approval(address,address,uint)")
	logTransferSigHash := crypto.Keccak256Hash(logTransferSig)
	logApprovalSigHash := crypto.Keccak256Hash(logApprovalSig)

	// 遍历事件日志
	for _, vLog := range logs {
		fmt.Printf("Log Block Number: %v\n", vLog.BlockNumber)
		fmt.Printf("Log Index: %v\n", vLog.Index)

		switch vLog.Topics[0].Hex() {
		case logTransferSigHash.Hex():
			fmt.Printf("Log Name: Transfer\n")

			// 事件日志类型转换
			var transferEvent LogTransfer
			err := contractAbi.UnpackIntoInterface(&transferEvent, "Transfer", vLog.Data)
			if err != nil {
				log.Fatal(err)
			}

			transferEvent.From = common.HexToAddress(vLog.Topics[1].Hex())
			transferEvent.To = common.HexToAddress(vLog.Topics[2].Hex())

			fmt.Printf("From: %s\n", transferEvent.From.Hex())
			fmt.Printf("To: %s\n", transferEvent.To.Hex())
			fmt.Printf("Tokens: %s\n", transferEvent.Tokens.String())
		case logApprovalSigHash.Hex():
			fmt.Printf("Log Name: Approval\n")

			var approveEvent LogApproval
			err := contractAbi.UnpackIntoInterface(&approveEvent, "Approval", vLog.Data)
			if err != nil {
				log.Fatal(err)
			}

			approveEvent.TokenOwner = common.HexToAddress(vLog.Topics[1].Hex())
			approveEvent.Spender = common.HexToAddress(vLog.Topics[2].Hex())

			fmt.Printf("Token Owner: %s\n", approveEvent.TokenOwner.Hex())
			fmt.Printf("Spender: %s\n", approveEvent.Spender.Hex())
			fmt.Printf("Tokens: %s\n", approveEvent.Tokens.String())
		}

		fmt.Printf("\n\n")
	}
}
