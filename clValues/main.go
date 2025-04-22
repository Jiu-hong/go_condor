package main

import (
	"context"
	"fmt"
	"go_condor/utils"
	"log"
	"math/big"
	"net/http"
	"time"

	"github.com/make-software/casper-go-sdk/v2/casper"
	"github.com/make-software/casper-go-sdk/v2/types"
	"github.com/make-software/casper-go-sdk/v2/types/clvalue"
	"github.com/make-software/casper-go-sdk/v2/types/clvalue/cltype"
	"github.com/make-software/casper-go-sdk/v2/types/key"
)

func main() {
	keys, err := casper.NewED25519PrivateKeyFromPEMFile(utils.KEYPATH)
	if err != nil {
		panic(err)
	}
	pubKey := keys.PublicKey()

	target, _ := casper.NewPublicKey("0106ed45915392c02b37136618372ac8dde8e0e3b8ee6190b2ca6db539b354ede4")

	mapvalue := clvalue.NewCLMap(cltype.String, cltype.Int32)
	mapvalue.Map.Append(*clvalue.NewCLString("ABC"), clvalue.NewCLInt32(10))
	mapvalue.Map.Append(*clvalue.NewCLString("XYZ"), clvalue.NewCLInt32(22000))
	mapvalue.Map.Append(*clvalue.NewCLString("DEF"), clvalue.NewCLInt32(10))
	mapvalue.Map.Append(*clvalue.NewCLString("DFIJ"), clvalue.NewCLInt32(22000))
	mapvalue.Map.Append(*clvalue.NewCLString("ABC"), clvalue.NewCLInt32(10))
	mapvalue.Map.Append(*clvalue.NewCLString("XYZ"), clvalue.NewCLInt32(22000))
	mapvalue.Map.Append(*clvalue.NewCLString("DEF"), clvalue.NewCLInt32(10))
	mapvalue.Map.Append(*clvalue.NewCLString("DFIJ"), clvalue.NewCLInt32(22000))

	listvalue := clvalue.NewCLList(cltype.String)
	listvalue.List.Append(*clvalue.NewCLString("ABC"))
	listvalue.List.Append(*clvalue.NewCLString("DEF"))

	resultcorrectvalue, _ := clvalue.NewCLResult(cltype.String, cltype.UInt32, *clvalue.NewCLString("ABC"), true)
	resulterrvalue, _ := clvalue.NewCLResult(cltype.Bool, cltype.UInt32, *clvalue.NewCLUInt32(10), false)

	optionresult := clvalue.NewCLOption(resulterrvalue)

	mapresult := clvalue.NewCLMap(cltype.String, &cltype.Option{Inner: &cltype.Result{}})
	mapresult.Map.Append(*clvalue.NewCLString("ABC"), optionresult)
	// mapresult := clvalue.NewCLOption(resulterrvalue)
	args := &types.Args{}
	args.AddArgument("string", *clvalue.NewCLString("Test")).
		AddArgument("u8", *clvalue.NewCLUint8(9)).
		AddArgument("u256", *clvalue.NewCLUInt256(big.NewInt(1_000_000_000_000_000))).
		AddArgument("option", clvalue.NewCLOption(*clvalue.NewCLUint8(9))).
		AddArgument("map", mapvalue). //NG
		AddArgument("list", listvalue).
		AddArgument("publickey", clvalue.NewCLPublicKey(target)).
		AddArgument("resultcorrectvalue", resultcorrectvalue).
		AddArgument("optionresult", optionresult).
		AddArgument("resulterr", resulterrvalue).
		AddArgument("mapresult", mapresult)

	entrypoint := "test2"
	packageHash, err := key.NewHash("40ad74eb43330f7fb496d6ea49df990e6583f51a01a7204a17a6217dbeb715d7")
	if err != nil {
		fmt.Println(err)
	}

	payload, err := types.NewTransactionV1Payload(
		types.InitiatorAddr{
			PublicKey: &pubKey,
		},
		types.Timestamp(time.Now().UTC()),
		1800000000000,
		utils.NETWORKNAME,
		types.PricingMode{
			Limited: &types.LimitedMode{
				PaymentAmount:     2500000000,
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
