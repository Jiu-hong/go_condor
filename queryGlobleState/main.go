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
	contractPackageHash := "hash-dc2530a47916c0fb042a94f084b51ae4ec04476933340169ad046b7c5ec5a078"
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
