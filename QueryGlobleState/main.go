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
	contractPackageHash := "hash-19cf434b80aa05d506f475a52da877240517a0ab238a49a54015e46e02649bbd"
	result, err := rpcClient.QueryLatestGlobalState(context.Background(), contractPackageHash, nil)
	if err != nil {
		fmt.Println(err)
	}

	b, err := json.MarshalIndent(result.StoredValue, "", "  ")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Print(string(b))
}
