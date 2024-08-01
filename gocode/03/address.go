package main

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
)

// 链接区块链钱包账户
func main() {
	// 链接钱包账户地址并获取钱包信息
	address := common.HexToAddress("0xf0b90BFb2e1418A37a7f38d0fEA92E2aa9e34389")
	fmt.Println(address.Hex())
	fmt.Println(address.MarshalText())
	fmt.Println(address.Bytes())
}
