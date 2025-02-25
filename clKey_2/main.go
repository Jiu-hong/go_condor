package main

import (
	"context"
	"fmt"
	"go_condor/utils"
	"log"
	"net/http"
	"time"

	"github.com/make-software/casper-go-sdk/v2/casper"
	"github.com/make-software/casper-go-sdk/v2/types"
	"github.com/make-software/casper-go-sdk/v2/types/clvalue"
	"github.com/make-software/casper-go-sdk/v2/types/key"
)

func main() {
	keys, err := casper.NewED25519PrivateKeyFromPEMFile(utils.KEYPATH)
	if err != nil {
		panic(err)
	}
	pubKey := keys.PublicKey()
	bid, _ := key.NewKey("bid-306633f962155a7d46658adb36143f28668f530454fe788c927cecf62e5964a1")
	unbond, _ := key.NewKey("unbond-2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a")
	withdraw, _ := key.NewKey("withdraw-2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a")
	bidaddr, _ := key.NewKey("bid-addr-00306633f962155a7d46658adb36143f28668f530454fe788c927cecf62e5964a1")
	dictionary, _ := key.NewKey("dictionary-2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a")
	balance, _ := key.NewKey("balance-2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a")
	chainspecregistry, _ := key.NewKey("chainspec-registry-0000000000000000000000000000000000000000000000000000000000000000")
	checksumregistry, _ := key.NewKey("checksum-registry-0000000000000000000000000000000000000000000000000000000000000000")
	balancehold, _ := key.NewKey("balance-hold-002a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a0000000000000000")
	args := &types.Args{}
	args.AddArgument("NewCLKeybid", clvalue.NewCLKey(bid)).
		AddArgument("NewCLKeyunbond", clvalue.NewCLKey(unbond)).
		AddArgument("NewCLKeywithdraw", clvalue.NewCLKey(withdraw)).
		AddArgument("NewCLKeybidaddr", clvalue.NewCLKey(bidaddr)).
		AddArgument("NewCLKeydictionary", clvalue.NewCLKey(dictionary)).
		AddArgument("NewCLKeybalance", clvalue.NewCLKey(balance)).
		AddArgument("NewCLKeychainspecregistry", clvalue.NewCLKey(chainspecregistry)).
		AddArgument("NewCLKeychecksumregistry", clvalue.NewCLKey(checksumregistry)).
		AddArgument("NewCLKeybalancehold", clvalue.NewCLKey(balancehold))

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

	rpcClient := casper.NewRPCClient(casper.NewRPCHandler(utils.ENDPOINT, http.DefaultClient))
	res, err := rpcClient.PutTransactionV1(context.Background(), *transaction)
	if err != nil {
		fmt.Println(err)
	}

	log.Println("TransactionV1 submitted:", res.TransactionHash.TransactionV1)

}
