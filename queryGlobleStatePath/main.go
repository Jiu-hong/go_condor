package main

import (
	"context"
	"encoding/json"
	"fmt"
	"go_condor/utils"
	"net/http"

	"github.com/make-software/casper-go-sdk/v2/casper"
)

func main() {

	rpcClient := casper.NewRPCClient(casper.NewRPCHandler(utils.ENDPOINT, http.DefaultClient))
	contractHash := "hash-082afb61e92f2b88013355502e17a2e3fa62a13eae2d134dc741c618b7eeb4d3"
	result, err := rpcClient.QueryLatestGlobalState(context.Background(), contractHash, []string{"nonce_executed"})
	if err != nil {
		fmt.Println(err)
	}

	b, err := json.MarshalIndent(result.StoredValue, "", "  ")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Print(string(b))
}
