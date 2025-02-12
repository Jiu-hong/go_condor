package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"os"
	"time"

	"go_condor/utils"

	"github.com/make-software/casper-go-sdk/v2/casper"
	"github.com/make-software/casper-go-sdk/v2/rpc"
	"github.com/make-software/casper-go-sdk/v2/types"
	"github.com/make-software/casper-go-sdk/v2/types/clvalue"
)

func main() {
	keys, err := casper.NewED25519PrivateKeyFromPEMFile(utils.KEYPATH)
	if err != nil {
		fmt.Println(err)
	}
	pubKey := keys.PublicKey()

	contractPath := "/home/ubuntu/mywork/mycontract_condor/contract/target/wasm32-unknown-unknown/release/contract.wasm"
	moduleBytes, err := os.ReadFile(contractPath)
	if err != nil {
		fmt.Println(err)
	}

	args := &types.Args{}
	args.AddArgument("name", *clvalue.NewCLString("Test")).
		AddArgument("symbol", *clvalue.NewCLString("test")).
		AddArgument("decimals", *clvalue.NewCLUint8(9)).
		AddArgument("total_supply", *clvalue.NewCLUInt256(big.NewInt(1_000_000_000_000_000))).
		AddArgument("events_mode", *clvalue.NewCLUint8(2)).
		AddArgument("enable_mint_burn", *clvalue.NewCLUint8(1))

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
				IsInstallUpgrade: true,
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

	rpcClient := rpc.NewClient(rpc.NewHttpHandler(utils.ENDPOINT, http.DefaultClient))
	res, err := rpcClient.PutTransactionV1(context.Background(), *transaction)
	if err != nil {
		fmt.Println(err)
	}

	log.Println("TransactionV1 submitted:", res.TransactionHash.TransactionV1)
}
