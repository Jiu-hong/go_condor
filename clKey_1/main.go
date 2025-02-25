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

	keyAccountHash, _ := key.NewKey("account-hash-bf06bdb1616050cea5862333d1f4787718f1011c95574ba92378419eefeeee59")
	keyHash, _ := key.NewKey("hash-8f6dec134948560997ba0236b54ce96679f5a51a35fe794df9c1e820538cbde3")

	keyUref, _ := key.NewKey("uref-7b12008bb757ee32caefb3f7a1f77d9f659ee7a4e21ad4950c4e0294000492eb-007")
	keyPackage, _ := key.NewKey("package-0000000000000000000000000000000000000000000000000000000000000000")
	keyTransfer, _ := key.NewKey("transfer-0404040404040404040404040404040404040404040404040404040404040404")
	keyDeploy, _ := key.NewKey("deploy-0505050505050505050505050505050505050505050505050505050505050505")
	eraId, _ := key.NewKey("era-42")

	system_entity_registry, _ := key.NewKey("system-contract-registry-0000000000000000000000000000000000000000000000000000000000000000")
	era_summary, _ := key.NewKey("era-summary-0000000000000000000000000000000000000000000000000000000000000000")
	args := &types.Args{}
	args.AddArgument("NewCLKeyAccountHash", clvalue.NewCLKey(keyAccountHash)).
		AddArgument("NewCKeyHash", clvalue.NewCLKey(keyHash)).
		AddArgument("NewCLKeyUref", clvalue.NewCLKey(keyUref)).
		AddArgument("NewCLKeyPackage", clvalue.NewCLKey(keyPackage)).
		AddArgument("NewCLKeyTransfer", clvalue.NewCLKey(keyTransfer)).
		AddArgument("NewCLKeyDeploy", clvalue.NewCLKey(keyDeploy)).
		AddArgument("NewCLKeyeraId", clvalue.NewCLKey(eraId)).
		AddArgument("NewCLKeysystem_entity_registry", clvalue.NewCLKey(system_entity_registry)).
		AddArgument("NewCLKeyera_summary", clvalue.NewCLKey(era_summary))

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
