// ok
package main

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"go_condor/utils"
	"math/big"
	"net/http"

	"github.com/make-software/casper-go-sdk/v2/casper"
	"github.com/make-software/casper-go-sdk/v2/types/clvalue"
)

func main() {
	keys, err := casper.NewED25519PrivateKeyFromPEMFile(utils.KEYPATH)
	if err != nil {
		panic(err)
	}
	accountPublicKey, err := casper.NewPublicKey("012488699f9a31e36ecf002675cd7186b48e6a735d10ec1b308587ca719937752c")
	if err != nil { return }
	amount := big.NewInt(100000000000)

	contractPath := "/home/ubuntu/mywork/mycontract_condor/contract/target/wasm32-unknown-unknown/release/contract.wasm"
	moduleBytes := utils.GetmoduleBytes(contractPath)
	session := casper.ExecutableDeployItem{
		ModuleBytes: &casper.ModuleBytes{
			ModuleBytes: hex.EncodeToString(moduleBytes),
			Args: (&casper.Args{}).
				AddArgument("target", clvalue.NewCLByteArray(accountPublicKey.AccountHash().Bytes())).
				AddArgument("amount", *clvalue.NewCLUInt512(amount)),
		},
	}

	payment := casper.StandardPayment(amount)

	deployHeader := casper.DefaultHeader()
	deployHeader.Account = keys.PublicKey()
	deployHeader.ChainName = utils.NETWORKNAME

	newDeploy, err := casper.MakeDeploy(deployHeader, payment, session)
	if err != nil {
		return
	}

	err = newDeploy.Sign(keys)
	if err != nil {
		panic(err)
	}

	handler := casper.NewRPCHandler(utils.ENDPOINT, http.DefaultClient)
	client := casper.NewRPCClient(handler)
	result, err := client.PutDeploy(context.Background(), *newDeploy)
	if err != nil {
		return
	}
	b, err := json.MarshalIndent(result, "", "  ")
    if err != nil {
        fmt.Println(err)
    }
    fmt.Print(string(b))
	// log.Println(result.DeployHash)
}

// contract-package-b0641a4c7ddc401b31950c69e0677da5ebc98938d2c7eaf081e398c82dcf7a72