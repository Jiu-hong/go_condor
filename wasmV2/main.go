package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"time"

	"go_condor/utils"

	"github.com/make-software/casper-go-sdk/v2/casper"
	"github.com/make-software/casper-go-sdk/v2/types"
	"github.com/make-software/casper-go-sdk/v2/types/clvalue"
	"github.com/make-software/casper-go-sdk/v2/types/clvalue/cltype"
	"github.com/make-software/casper-go-sdk/v2/types/key"
)

func main() {
	keys, err := casper.NewED25519PrivateKeyFromPEMFile(utils.KEYPATH)
	if err != nil {
		fmt.Println(err)
	}
	pubKey := keys.PublicKey()

	// contractPath := "/home/ubuntu/mywork/mycontract_condor/contract/target/wasm32-unknown-unknown/release/contract.wasm"
	// moduleBytes, err := os.ReadFile(contractPath)
	hex_moduleBytes := utils.CONTRACT_HEX_CODE
	moduleBytes, err := hex.DecodeString(hex_moduleBytes)
	if err != nil {
		fmt.Println(err)
	}

	associatedAccount1, err := key.NewAccountHash("13f3cbf02c716f397239db88180fca7704cc1703f77a4ee925f719b82174a049")
	if err != nil {
		fmt.Println(err)
	}
	associatedAccount2, err := key.NewAccountHash("b34999881cd9bb0ad66378a527c6a8ad0590ca98e53f57e103090501c4406aa8")
	if err != nil {
		fmt.Println(err)
	}
	associatedAccount3, err := key.NewAccountHash("4df9c55556635622093cf62cb92b4645b6f95cef7ea163dc3d3dbcb3a1a2291f")
	if err != nil {
		fmt.Println(err)
	}
	accountsList := clvalue.NewCLList(cltype.NewByteArray(32))
	accountsList.List.Append(clvalue.NewCLByteArray(associatedAccount1.Hash[:]))
	accountsList.List.Append(clvalue.NewCLByteArray(associatedAccount2.Hash[:]))
	accountsList.List.Append(clvalue.NewCLByteArray(associatedAccount3.Hash[:]))

	listvalue := clvalue.NewCLList(cltype.UInt8)
	listvalue.List.Append(*clvalue.NewCLUint8(1))
	listvalue.List.Append(*clvalue.NewCLUint8(1))
	listvalue.List.Append(*clvalue.NewCLUint8(1))

	args := &types.Args{}
	args.AddArgument("action", *clvalue.NewCLString("set_all")).
		AddArgument("deployment_thereshold", *clvalue.NewCLUint8(2)).
		AddArgument("key_management_threshold", *clvalue.NewCLUint8(9)).
		AddArgument("accounts", accountsList).
		AddArgument("weights", listvalue).
		AddArgument("deploy_type", *clvalue.NewCLString("WalletInitialization")).
		AddArgument("owner_0", *clvalue.NewCLString("01366d77126a722e1adce6ad0bf2fbdbcdc573eb5c9c338a1097f1837f3dc4ef88")).
		AddArgument("owner_1", *clvalue.NewCLString("01a36e511344bcecbbe8082bec459c79a66d777d75d33bbed88e7fbd242e33f65d")).
		AddArgument("owner_2", *clvalue.NewCLString("01c176c9c1281eba037efbbe8c46dc57dcf49a22b23de9cafe9477998993e31972"))

	payload, err := types.NewTransactionV1Payload(
		types.InitiatorAddr{
			PublicKey: &pubKey,
		},
		types.Timestamp(time.Now().UTC()),
		1800000000000,
		utils.NETWORKNAME,
		types.PricingMode{
			Limited: &types.LimitedMode{
				PaymentAmount:     10000000000,
				GasPriceTolerance: 1,
				StandardPayment:   true,
			},
		},
		types.NewNamedArgs(args),
		types.TransactionTarget{
			Session: &types.SessionTarget{
				ModuleBytes:      moduleBytes,
				Runtime:          types.NewVmCasperV1TransactionRuntime(),
				IsInstallUpgrade: false,
			},
		},
		types.TransactionEntryPoint{
			Call: &struct{}{},
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
