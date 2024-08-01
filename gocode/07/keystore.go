package main

import (
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"io/ioutil"
	"log"
	"os"
)

// 通过新创建方式生成加密钱包私钥
func createKs() {
	// 初始化钱包私钥存储地址
	ks := keystore.NewKeyStore("./tmp", keystore.StandardScryptN, keystore.StandardScryptP)
	password := "secret"
	// 创建新的钱包
	account, err := ks.NewAccount(password)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(account.Address.Hex())
}

// 通过导入方式生成加密钱包私钥
func importKs() {
	// 旧的钱包私钥存储地址
	file := ""
	// 初始化信息新的钱包私钥存储地址
	ks := keystore.NewKeyStore("./tmp", keystore.StandardScryptN, keystore.StandardScryptP)
	jsonBytes, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}
	// 加密私钥口令
	password := "secret"
	// 新的加密口令
	newPassword := "secretNew"
	account, err := ks.Import(jsonBytes, password, newPassword)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(account.Address.Hex())

	if err := os.Remove(file); err != nil {
		log.Fatal(err)
	}
}

// 生成加密钱包私钥
func main() {
	createKs()
	//importKs()

}
