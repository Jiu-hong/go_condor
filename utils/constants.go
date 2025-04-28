package utils

import "os"

const NETWORKNAME = "integration-test"
const ENDPOINT = "http://node.integration.casper.network:7777/rpc"

// const NETWORKNAME = "casper-jiuhong-test-jh-1"
// const ENDPOINT = "http://35.87.247.87:7777/rpc"

// const NETWORKNAME = "dev-net"
// const ENDPOINT = "http://54.174.173.4:7777/rpc"
const SIGNEDTRANSACTIONPATH = "/home/ubuntu/keys/test0/secret_key.pem"
const KEYPATH = "/home/ubuntu/keys/test0/secret_key.pem"
const TTL = 180000000000

func GetmoduleBytes(contractPath string) []byte {
	moduleBytes, err := os.ReadFile(contractPath)
	if err != nil {
		panic(err)
	}

	return moduleBytes
}
