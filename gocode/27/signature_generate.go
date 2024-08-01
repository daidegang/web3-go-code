package main

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"log"
)

// 生成签名
func main() {
	// 加载私钥
	privateKey, err := crypto.HexToECDSA("ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80")
	if err != nil {
		log.Fatal(err)
	}

	// 初始化需要签名数据
	data := []byte("hello")
	hash := crypto.Keccak256Hash(data)
	fmt.Printf("%x\n", hash.Hex())

	// 使用私钥对数据进行签名
	signature, err := crypto.Sign(hash.Bytes(), privateKey)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%x\n", hexutil.Encode(signature))
}
