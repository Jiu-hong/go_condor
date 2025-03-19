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

	associatedAccount, err := key.NewAccountHash("80cef79f0451fdfb21084aaab8a4811e27dfc262c970d77675e3cad5394ef1f7")
	if err != nil {
		fmt.Println(err)
	}

	accountsList := clvalue.NewCLList(cltype.NewByteArray(32))
	accountsList.List.Append(clvalue.NewCLByteArray(associatedAccount.Hash[:]))

	listvalue := clvalue.NewCLList(cltype.UInt8)
	listvalue.List.Append(*clvalue.NewCLUint8(1))

	args := &types.Args{}
	args.AddArgument("action", *clvalue.NewCLString("set_all")).
		AddArgument("deployment_thereshold", *clvalue.NewCLUint8(2)).
		AddArgument("key_management_threshold", *clvalue.NewCLUint8(9)).
		AddArgument("accounts", accountsList).
		AddArgument("weights", listvalue).
		AddArgument("deploy_type", *clvalue.NewCLString("WalletInitialization")).
		AddArgument("owner_0", *clvalue.NewCLString("0203c34ddd4dcddfd5c0082cadf24613597712eb92230b901089f469170b44a569a8")).
		AddArgument("owner_1", *clvalue.NewCLString("0203c34ddd4dcddfd5c0082cadf24613597712eb92230b901089f469170b44a569a8")).
		AddArgument("owner_2", *clvalue.NewCLString("0203c34ddd4dcddfd5c0082cadf24613597712eb92230b901089f469170b44a569a8"))

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
