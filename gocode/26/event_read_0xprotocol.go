package main

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math/big"
	"strconv"
	"strings"
	exchange "web3-go-code/contract/bindings"
)

// 定义事件日志签名类型匹配结构体
type LogFill struct {
	Maker                  common.Address
	Taker                  common.Address
	FeeRecipient           common.Address
	MakerToken             common.Address
	TakerToken             common.Address
	FilledMakeTokenAmount  *big.Int
	FilledTakerTokenAmount *big.Int
	PaidMakerFee           *big.Int
	PaidTakerFee           *big.Int
	Tokens                 [32]byte
	OrderHash              [32]byte
}

type LogCancel struct {
	Maker                     common.Address
	FeeRecipient              common.Address
	MakerToken                common.Address
	TakerToken                common.Address
	CancelledMakerTokenAmount *big.Int
	CancelledTakerTokenAmount *big.Int
	Tokens                    [32]byte
	OrderHash                 [32]byte
}

type LogError struct {
	ErrorID   uint8
	OrderHash [32]byte
}

// 读取0x Protocol事件日志
func main() {
	// 初始化以太坊客户端
	client, err := ethclient.Dial("http://127.0.0.1:8545")
	if err != nil {
		log.Fatal(err)
	}

	// 加载合约地址
	contractAddress := common.HexToAddress("0x5fc8d32690cc91d4c39d9d3abcbd16989f875707")

	// 初始化查询区块信息
	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(0),
		ToBlock:   big.NewInt(6383488),
		Addresses: []common.Address{
			contractAddress,
		},
	}

	// 查询日志
	logs, err := client.FilterLogs(context.Background(), query)
	if err != nil {
		log.Fatal(err)
	}

	// 解析JSON abi
	contractAbi, err := abi.JSON(strings.NewReader(string(exchange.ExchangeABI)))
	if err != nil {
		log.Fatal(err)
	}

	// 生成事件日志函数签名keccak256摘要
	logFillSig := []byte("LogFill(address,address,address,address,address,uint256,uint256,uint256,uint256,bytes32,bytes32)")
	logCancelSig := []byte("logCancel(address,address,address,address,uint256,uint256,bytes32,bytes32)")
	logErrorSig := []byte("logError(uint8,bytes32)")
	logFillSigHash := crypto.Keccak256Hash(logFillSig)
	logCancelSigHash := crypto.Keccak256Hash(logCancelSig)
	logErrorSigHash := crypto.Keccak256Hash(logErrorSig)

	logFillEvent := common.HexToHash(logFillSigHash.String())
	logCancelEvent := common.HexToHash(logCancelSigHash.String())
	logErrorEvent := common.HexToHash(logErrorSigHash.String())

	// 迭代所有日志并设置一个switch语句来按事件日志类型过滤
	for _, vLog := range logs {
		fmt.Printf("Log Block Number: %d\n", vLog.BlockNumber)
		fmt.Printf("Log Index: %d\n", vLog.Index)

		switch vLog.Topics[0].Hex() {
		case logFillEvent.Hex():
			fmt.Printf("Log Name: LogFill\n")

			var fillevent LogFill

			err := contractAbi.UnpackIntoInterface(&fillevent, "LogFill", vLog.Data)
			if err != nil {
				log.Fatal(err)
			}

			fillevent.Maker = common.HexToAddress(vLog.Topics[1].Hex())
			fillevent.FeeRecipient = common.HexToAddress(vLog.Topics[2].Hex())
			fillevent.Tokens = vLog.Topics[3]

			fmt.Printf("Maker: %s\n", fillevent.Maker.Hex())
			fmt.Printf("Taker: %s\n", fillevent.Taker.Hex())
			fmt.Printf("Fee Recipient: %s\n", fillevent.FeeRecipient.Hex())
			fmt.Printf("Maker Token: %s\n", fillevent.MakerToken.Hex())
			fmt.Printf("Taker Token: %s\n", fillevent.TakerToken.Hex())
			fmt.Printf("Filled Maker Token Amount: %s\n", fillevent.FilledMakeTokenAmount.String())
			fmt.Printf("Filled Taker Token Amount: %s\n", fillevent.FilledTakerTokenAmount.String())
			fmt.Printf("Paid Maker Fee: %s\n", fillevent.PaidMakerFee.String())
			fmt.Printf("Paid Taker Fee: %s\n", fillevent.PaidTakerFee.String())
			fmt.Printf("Tokens: %s\n", hexutil.Encode(fillevent.Tokens[:]))
			fmt.Printf("Order Hash: %s\n", hexutil.Encode(fillevent.OrderHash[:]))
		case logCancelEvent.Hex():
			fmt.Printf("Log Name: LogCancel\n")

			var cancelEvent LogCancel

			err := contractAbi.UnpackIntoInterface(&cancelEvent, "LogCancel", vLog.Data)
			if err != nil {
				log.Fatal(err)
			}

			cancelEvent.Maker = common.HexToAddress(vLog.Topics[1].Hex())
			cancelEvent.FeeRecipient = common.HexToAddress(vLog.Topics[2].Hex())
			cancelEvent.Tokens = vLog.Topics[3]

			fmt.Printf("Maker: %s\n", cancelEvent.Maker.Hex())
			fmt.Printf("Fee Recipient: %s\n", cancelEvent.FeeRecipient.Hex())
			fmt.Printf("Maker Token: %s\n", cancelEvent.MakerToken.Hex())
			fmt.Printf("Taker Token: %s\n", cancelEvent.TakerToken.Hex())
			fmt.Printf("Cancelled Maker Token Amount: %s\n", cancelEvent.CancelledMakerTokenAmount.String())
			fmt.Printf("Canceled Taker Token Amount: %s\n", cancelEvent.CancelledTakerTokenAmount.String())
			fmt.Printf("Tokens: %s\n", hexutil.Encode(cancelEvent.Tokens[:]))
			fmt.Printf("Order Hash: %s\n", hexutil.Encode(cancelEvent.OrderHash[:]))

		case logErrorEvent.Hex():
			fmt.Printf("Log Name: LogError\n")
			errorID, err := strconv.ParseInt(vLog.Topics[1].Hex(), 16, 64)
			if err != nil {
				log.Fatal(err)
			}

			errorEvent := &LogError{
				ErrorID:   uint8(errorID),
				OrderHash: vLog.Topics[2],
			}
			fmt.Printf("Error ID: %d\n", errorEvent.ErrorID)
			fmt.Printf("Order Hash: %s\n", hexutil.Encode(errorEvent.OrderHash[:]))
		}

		fmt.Printf("\n\n")
	}
}
