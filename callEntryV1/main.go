// ok
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"go_condor/utils"
	"math/big"
	"net/http"

	"github.com/make-software/casper-go-sdk/v2/casper"
	"github.com/make-software/casper-go-sdk/v2/types/clvalue"
	"github.com/make-software/casper-go-sdk/v2/types/key"
)

func main() {
	keys, err := casper.NewED25519PrivateKeyFromPEMFile(utils.KEYPATH)
	if err != nil {
		panic(err)
	}

	deployHeader := casper.DefaultHeader()
	deployHeader.Account = keys.PublicKey()
	deployHeader.ChainName = utils.NETWORKNAME

	// header.Timestamp = types.Timestamp(time.Now())
	// payment := types.StandardPayment(big.NewInt(4000000000))
	payment := casper.StandardPayment(big.NewInt(4000000000))

	sessionArgs := &casper.Args{}
	key1, err := key.NewKey("account-hash-bf06bdb1616050cea5862333d1f4787718f1011c95574ba92378419eefeeee59")
	if err != nil {
		panic(err)
	}
	sessionArgs.AddArgument("amount", *clvalue.NewCLUInt256(big.NewInt(2500000000))).
		AddArgument("owner", clvalue.NewCLKey(key1))

	contractHash, err := key.NewContract("b0641a4c7ddc401b31950c69e0677da5ebc98938d2c7eaf081e398c82dcf7a72")
	if err != nil {
		panic(err)
	}
	varVal := json.Number("1")
	session := casper.ExecutableDeployItem{
		StoredVersionedContractByHash: &casper.StoredVersionedContractByHash{
			Hash:       contractHash,
			EntryPoint: "apple",
			Version:    &varVal,
			Args:       sessionArgs,
		},
	}

	deploy, err := casper.MakeDeploy(deployHeader, payment, session)
	if err != nil {
		panic(err)
	}
	err = deploy.Sign(keys)
	if err != nil {
		panic(err)
	}

	// rpcClient := rpc.NewClient(rpc.NewHttpHandler(utils.ENDPOINT, http.DefaultClient))
	// _, err = rpcClient.PutDeploy(context.Background(), *deploy)
	// if err != nil {
	// 	panic(err)
	// }
	handler := casper.NewRPCHandler(utils.ENDPOINT, http.DefaultClient)
	client := casper.NewRPCClient(handler)
	result, err := client.PutDeploy(context.Background(), *deploy)
	if err != nil {
		return
	}
	b, err := json.MarshalIndent(result, "", "  ")
    if err != nil {
        fmt.Println(err)
    }
    fmt.Print(string(b))
}
