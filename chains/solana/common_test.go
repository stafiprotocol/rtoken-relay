package solana_test

import (
	"testing"

	"github.com/stafiprotocol/rtoken-relay/chains/solana"
	solCommon "github.com/tpkeeper/solana-go-sdk/common"
)

func TestGetMultisigTxAccountPubkey(t *testing.T) {

	seed, pubkey := solana.GetMultisigTxAccountPubkey(
		solCommon.PublicKeyFromString("Gnr9LuHUh85Dt7Qr3tayXrxFAEn32jRDfsgTAyywFhyh"),
		solana.MultisigTxStakeType,
		119)
	t.Log(seed, pubkey.ToBase58())
}
