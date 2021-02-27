module github.com/stafiprotocol/rtoken-relay

go 1.13

require (
	github.com/ChainSafe/log15 v1.0.0
	github.com/cosmos/cosmos-sdk v0.41.3
	github.com/ethereum/go-ethereum v1.9.25
	github.com/itering/scale.go v0.7.0
	github.com/itering/substrate-api-rpc v0.2.0
	github.com/stafiprotocol/chainbridge v0.0.0-20210122054647-25195c4be148
	github.com/stafiprotocol/go-substrate-rpc-client v1.0.2
	github.com/stretchr/testify v1.7.0
	github.com/urfave/cli/v2 v2.3.0
	golang.org/x/crypto v0.0.0-20201221181555-eec23a3978ad
	gotest.tools v2.2.0+incompatible
)

replace github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1
