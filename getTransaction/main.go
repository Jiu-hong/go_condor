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
	transactionHash := "4de6396bb92a04e9106296bf6ca45283917d6d8c2cbc72569ee789aa31ba3034"
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