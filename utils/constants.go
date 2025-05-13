package utils

import "os"

// const NETWORKNAME = "casper-test"
// const ENDPOINT = "https://node.testnet.casper.network/rpc"

// const NETWORKNAME = "integration-test"
// const ENDPOINT = "http://node.integration.casper.network:7777/rpc"

const NETWORKNAME = "casper"
const ENDPOINT = "http://34.218.248.90:7777/rpc"

// const NETWORKNAME = "casper-jiuhong-test-jh-1"
// const ENDPOINT = "http://35.87.247.87:7777/rpc"

// const NETWORKNAME = "dev-net"
// const ENDPOINT = "http://54.174.173.4:7777/rpc"
const SIGNEDTRANSACTIONPATH = "signed_transaction"
const KEYPATH = "/home/ubuntu/keys/test0/secret_key.pem"
const TTL = 180000000000

func GetmoduleBytes(contractPath string) []byte {
	moduleBytes, err := os.ReadFile(contractPath)
	if err != nil {
		panic(err)
	}

	return moduleBytes
}
