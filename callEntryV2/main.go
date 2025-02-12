// NG
package main

import (
	"context"
	"fmt"
	"go_condor/utils"
	"log"
	"net/http"
	"time"

	"github.com/make-software/casper-go-sdk/v2/casper"
	"github.com/make-software/casper-go-sdk/v2/rpc"
	"github.com/make-software/casper-go-sdk/v2/types"
	"github.com/make-software/casper-go-sdk/v2/types/key"
)

func main() {
	keys, err := casper.NewED25519PrivateKeyFromPEMFile(utils.KEYPATH)
	if err != nil {
		panic(err)
	}
	pubKey := keys.PublicKey()

	args := &types.Args{}

	entrypoint := "apple"
	// hash, err := key.NewHash("a5542d422cc7102165bde32f8c8aa460a81dc64105b03efbcd9c612a7721dadb")
	packageHash, err := key.NewHash("e48c5b9631c3a2063e61826d6e52181ea5d6fe35566bf994134caa26fce16586")
	if err != nil {
		fmt.Println(err)
	}
	// name := "my_hash"

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
			// // by contract hash
			// Stored: &types.StoredTarget{
			// 	ID: types.TransactionInvocationTarget{
			// 		// ByHash: &hash,
			// 		ByName: &name,
			// 	},
			// 	Runtime: types.NewVmCasperV1TransactionRuntime(),
			// },
			// // by package name
			// Stored: &types.StoredTarget{
			// 	ID: types.TransactionInvocationTarget{
			// 		ByPackageName: &types.ByPackageNameInvocationTarget{Name: name},
			// 	},
			// 	Runtime: types.NewVmCasperV1TransactionRuntime(),
			// },
			// by package hash
			Stored: &types.StoredTarget{
				ID: types.TransactionInvocationTarget{
					ByPackageHash: &types.ByPackageHashInvocationTarget{Addr: packageHash},
				},
				Runtime: types.NewVmCasperV1TransactionRuntime(),
			},
		},
		types.TransactionEntryPoint{
			Custom: &entrypoint,
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

}
