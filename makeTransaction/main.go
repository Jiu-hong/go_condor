// ok
package main

import (
	"encoding/json"
	"fmt"
	"go_condor/utils"
	"math/big"
	"time"

	"github.com/make-software/casper-go-sdk/v2/casper"
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
				PaymentAmount:     100000000,
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

	b, err := json.MarshalIndent(transaction, "", "  ")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Print(string(b))

}
