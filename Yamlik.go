package main

import (
	"awesomeProject/config"
	"context"
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"gitlab.com/distributed_lab/kit/kv"
	"log"
	"math/big"
	"time"

	/*"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient" */

)



func MakeTransaction(){
	cfg := config.NewConfig(kv.MustFromEnv())

	eth := cfg.EthClient()

	privateKey, err := crypto.HexToECDSA("fad9c8855b740a0b7ed4c221dbad0f33a83a49cad6b3fe8d5817ac83d38b6a19")
	if err != nil {
		log.Fatal(err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)



	nonce, err := eth.PendingNonceAt(context.Background(), fromAddress )
	if err != nil {
		log.Fatal(err)
	}
	value := big.NewInt(1000000000000000000)
	gasLimit := uint64(21000)
	gasPrice, err := eth.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	toAddress := common.HexToAddress("0x4592d8f8d7b001e72cb26a73e4fa1806a51ac79d")
	var data []byte
	tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, data)
	chainID, err := eth.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Fatal(err)
	}

	err = eth.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("tx sent: %s", signedTx.Hash().Hex())
	fmt.Println("Transaction was made successfully")

}

func main() {
	d := time.NewTicker(3 * time.Second)
	MyChannel := make(chan bool)

	go func() {
		time.Sleep(9 * time.Second)
		MyChannel <- true
	}()

	for {
		select {
		case <-MyChannel:
			fmt.Println("Completed!")
			return
		case tm := <-d.C:
			MakeTransaction()
			fmt.Println("The Current time is: ", tm)
		}
	}

}
