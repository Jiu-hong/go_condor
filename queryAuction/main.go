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

	result, err := rpcClient.GetAuctionInfoV1ByHeight(context.Background(), 4793042)
	if err != nil {
		fmt.Println(err)
	}

	b, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Print(string(b))
}
