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

	// query contract_hash
	rpcClient := casper.NewRPCClient(casper.NewRPCHandler(utils.ENDPOINT, http.DefaultClient))
	contractHash := "hash-4de2ec67e8e82959a049df08cc76af3034c7347ba414334d5ec45d0e3e9e3c5f"
	result, err := rpcClient.QueryLatestGlobalState(context.Background(), contractHash, nil)
	if err != nil {
		fmt.Println(err)
	}
	// get name_keys
	namedKeys := result.StoredValue.Contract.NamedKeys

	token_owners_uref := ""
	// get the uref under name_keys
	for _, namedKey := range namedKeys {
		if namedKey.Name == "token_owners" {
			token_owners_uref = namedKey.Key.ToPrefixedString()
		}
	}
	b, err := json.MarshalIndent(namedKeys, "", "  ")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(b))
	fmt.Println("token_owners_uref:", token_owners_uref)

	// get the dictionary_value
	stateRootHash := "e057c05a5db2a13a27ce831d573f977fa4fceb21dbfb51b58e8a6ef7c3d41c6a"
	result1, err := rpcClient.GetDictionaryItem(context.Background(), &stateRootHash, token_owners_uref, "1")
	if err != nil {
		fmt.Println(err)
	}
	value, err := result1.StoredValue.CLValue.Value()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("value:", value)

}
