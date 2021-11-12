package cosmos_test

import (
	"encoding/hex"
	// substrateTypes "github.com/stafiprotocol/go-substrate-rpc-client/types"
	"github.com/stafiprotocol/rtoken-relay/chains/cosmos"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetBondUnBondProposalId(t *testing.T) {
	// bts := cosmos.GetBondUnBondProposalId(substrateTypes.NewHash([]byte{2, 2}), substrateTypes.NewU128(*big.NewInt(1)), substrateTypes.NewU128(*big.NewInt(10)), 5)
	// t.Log(hex.EncodeToString(bts))

	bts, err := hex.DecodeString("cbec7f11e2cc6793c0ff2419b6be51d9fe41dbf121bdfe1a8efea08a6956a71200000000000000000000000000000000724ae66f00000000000000000000000000000000000000d4")
	assert.NoError(t, err)
	id, bond, unbond, seq, err := cosmos.ParseBondUnBondProposalId(bts)
	assert.NoError(t, err)
	t.Log(id, bond, unbond, seq)
}
