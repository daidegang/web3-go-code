package main

import (
	"bytes"
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"log"
)

// 验证签名
func main() {
	// 加载私钥
	privateKey, err := crypto.HexToECDSA("ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80")
	if err != nil {
		log.Fatal(err)
	}

	// 私钥生成公钥
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}

	// 生成字节格式的公钥
	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)

	// 生成原始数据哈希
	data := []byte("hello world")
	hash := crypto.Keccak256Hash(data)
	fmt.Printf("%x\n", hash.Hex())

	// 使用私钥对数据进行签名
	signature, err := crypto.Sign(hash.Bytes(), privateKey)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf(hexutil.Encode(signature))

	// 调用椭圆曲线签名方法恢复检索签名者的公钥
	sigPublicKey, err := crypto.Ecrecover(hash.Bytes(), signature)
	if err != nil {
		log.Fatal(err)
	}

	// 验证公钥字节信息是否匹配(验证方式一)
	matches := bytes.Equal(sigPublicKey, publicKeyBytes)
	fmt.Println(matches)

	// 返回ECDSA类型中的签名公钥(验证方式二)
	signPublicECDSA, err := crypto.SigToPub(hash.Bytes(), signature)
	if err != nil {
		log.Fatal(err)
	}
	signPublicBytes := crypto.FromECDSAPub(signPublicECDSA)
	matches = bytes.Equal(signPublicBytes, publicKeyBytes)
	fmt.Println(matches)

	// 如果公钥与签名的签名者匹配（验证方式三）
	signatureNoRecoverID := signature[:len(signature)-1] // remove recovery ID
	verified := crypto.VerifySignature(publicKeyBytes, hash.Bytes(), signatureNoRecoverID)
	fmt.Println(verified)
}
