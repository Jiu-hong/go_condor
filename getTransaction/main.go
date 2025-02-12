// ok
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/make-software/casper-go-sdk/v2/casper"
)

func main() {
	handler := casper.NewRPCHandler("http://node.integration.casper.network:7777/rpc", http.DefaultClient)
	client := casper.NewRPCClient(handler)
	transactionHash := "ad4e321b81f7deefe4f6aefa25c989ee2e05c86ff5c87795eace0d5d3772417d"
	deploy, err := client.GetTransactionByTransactionHash(context.Background(), transactionHash)
	if err != nil {
		return
	}
	b, err := json.MarshalIndent(deploy, "", "  ")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Print(string(b))
}
