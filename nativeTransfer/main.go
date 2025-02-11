// ok
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"go_condor/utils"
	"log"
	"math/big"
	"net/http"
	"time"

	"github.com/make-software/casper-go-sdk/v2/casper"
	"github.com/make-software/casper-go-sdk/v2/rpc"
	"github.com/make-software/casper-go-sdk/v2/types"
	"github.com/make-software/casper-go-sdk/v2/types/clvalue"
)

func main() {
	keys, err := casper.NewED25519PrivateKeyFromPEMFile(utils.KEYPATH)
	if err != nil {
		panic(err)
	}
	pubKey := keys.PublicKey()

	target, err := casper.NewPublicKey("0106ed45915392c02b37136618372ac8dde8e0e3b8ee6190b2ca6db539b354ede4")
	if err != nil {
		panic(err)
	}

	args := &types.Args{}
	args.AddArgument("target", clvalue.NewCLPublicKey(target)).
		AddArgument("amount", *clvalue.NewCLUInt512(big.NewInt(2500000000)))

	payload, err := types.NewTransactionV1Payload(
		types.InitiatorAddr{
			PublicKey: &pubKey,
		},
		types.Timestamp(time.Now().UTC()),
		1800000000000,
		utils.NETWORKNAME,
		types.PricingMode{
			Limited: &types.LimitedMode{
				PaymentAmount:     100000,
				GasPriceTolerance: 1,
				StandardPayment:   true,
			},
		},
		types.NewNamedArgs(args),
		types.TransactionTarget{
			Native: &struct{}{},
		},
		types.TransactionEntryPoint{
			Transfer: &struct{}{},
		},
		types.TransactionScheduling{
			Standard: &struct{}{},
		},
	)
	if err != nil {
		fmt.Println(err)
	}
	transaction, err := types.MakeTransactionV1(payload)
	if err != nil {
		fmt.Println(err)
	}

	err = transaction.Sign(keys)
	if err != nil {
		fmt.Println(err)
	}

	rpcClient := rpc.NewClient(rpc.NewHttpHandler(utils.ENDPOINT, http.DefaultClient))
	res, err := rpcClient.PutTransactionV1(context.Background(), *transaction)
	if err != nil {
		fmt.Println(err)
	}

	log.Println("TransactionV1 submitted:", res.TransactionHash.TransactionV1)

	time.Sleep(time.Second * 10)
	transactionRes, err := rpcClient.GetTransactionByTransactionHash(context.Background(), res.TransactionHash.TransactionV1.ToHex())
	if err != nil {
		fmt.Println(err)
	}
	// fmt.Println(transactionRes)
	b, err := json.MarshalIndent(transactionRes, "", "  ")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Print(string(b))

}
