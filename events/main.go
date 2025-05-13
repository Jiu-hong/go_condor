package main

import (
	"context"
	"fmt"
	"go_condor/utils"
	"net/http"

	"github.com/make-software/casper-go-sdk/v2/casper"
	"github.com/make-software/ces-go-parser/v2"
)

func main() {
	client := casper.NewRPCClient(casper.NewRPCHandler(utils.ENDPOINT, http.DefaultClient))
	// the transaction hash installing contract
	transactionInfo, err := client.GetTransactionByTransactionHash(context.Background(), "5e45de4d3a28d073a2b27531d32a0c12b0697290fcd9551a99066180e38b2649")
	if err != nil {
		panic(err)
	}
	executionResult := transactionInfo.ExecutionInfo.ExecutionResult

	// contract hash
	contracthash, err := casper.NewHash("0c30389226f0c938d69230f6bed4fbb3e54479910ec9c917a33fba7f342eb6c5")
	if err != nil {
		panic(err)
	}
	cesParser, err := ces.NewParser(client, []casper.Hash{contracthash})
	if err != nil {
		panic(err)
	}

	parseResult, err := cesParser.ParseExecutionResults(*executionResult)

	if err != nil {
		panic(err)
	}
	for _, result := range parseResult {
		if result.Error == nil {
			eventName := result.Event.Name
			eventData := result.Event.Data
			fmt.Println(eventName, eventData)
		} else {
			fmt.Println(result.Error)
		}
	}
}
