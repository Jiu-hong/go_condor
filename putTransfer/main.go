// ok
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"go_condor/utils"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/make-software/casper-go-sdk/v2/casper"
)

func main() {
	transaction_bytes, err := os.ReadFile(utils.SIGNEDTRANSACTIONPATH)
	if err != nil {
		fmt.Println(err)
	}

	var transaction casper.TransactionV1
	err = json.Unmarshal(transaction_bytes, &transaction)
	if err != nil {
		fmt.Println("Failed to unmarshal transaction:", err)
		return
	}

	rpcClient := casper.NewRPCClient(casper.NewRPCHandler(utils.ENDPOINT, http.DefaultClient))

	res, err := rpcClient.PutTransactionV1(context.Background(), transaction)
	if err != nil {
		fmt.Println(err)
	}

	log.Println("TransactionV1 submitted:", res.TransactionHash.TransactionV1)

	time.Sleep(time.Second * 10)
	transactionRes, err := rpcClient.GetTransactionByTransactionHash(context.Background(), res.TransactionHash.TransactionV1.ToHex())
	if err != nil {
		fmt.Println(err)
	}

	b, err := json.MarshalIndent(transactionRes, "", "  ")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Print(string(b))

}
