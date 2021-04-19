package cosmos_test

import (
	"encoding/hex"
	substrateTypes "github.com/stafiprotocol/go-substrate-rpc-client/types"
	"github.com/stafiprotocol/rtoken-relay/chains/cosmos"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

func TestGetBondUnBondProposalId(t *testing.T) {
	bts := cosmos.GetBondUnBondProposalId(substrateTypes.NewHash([]byte{2, 2}), substrateTypes.NewU128(*big.NewInt(1)), substrateTypes.NewU128(*big.NewInt(10)), 5)
	t.Log(hex.EncodeToString(bts))
	id, bond, unbond, seq, err := cosmos.ParseBondUnBondProposalId(bts)
	assert.NoError(t, err)
	t.Log(id, bond, unbond, seq)
}
