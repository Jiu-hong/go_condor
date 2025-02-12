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
	transactionHash := "f9fa4d7c26d77758ea0ef70184ae9210c1c915721bb541d7e9dc4e2d09c9954b"
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
