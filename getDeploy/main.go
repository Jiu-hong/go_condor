// ok
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/make-software/casper-go-sdk/v2/casper"
)

func main() {
	handler := casper.NewRPCHandler("http://node.integration.casper.network:7777/rpc", http.DefaultClient)
	client := casper.NewRPCClient(handler)
	deployHash := "fe0d150c43b8093492043c24d700122443c1b4745d534a5df25304501b96b7b2"
	deploy, err := client.GetDeploy(context.Background(), deployHash)
	if err != nil {
		return
	}
	b, err := json.MarshalIndent(deploy, "", "  ")
    if err != nil {
        fmt.Println(err)
    }
    fmt.Print(string(b))
}