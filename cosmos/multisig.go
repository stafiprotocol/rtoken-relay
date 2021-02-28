package cosmos

import (
	kMultiSig "github.com/cosmos/cosmos-sdk/crypto/keys/multisig"
	"github.com/cosmos/cosmos-sdk/crypto/types/multisig"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
)

func (c *Client) CreateMultiSigTx(rawTx []byte, signatures [][]byte) (txBts []byte, err error) {
	accountMultiSign, err := c.clientCtx.AccountRetriever.GetAccount(c.clientCtx, c.clientCtx.GetFromAddress())
	if err != nil {
		return nil, err
	}
	legacyPubkey := accountMultiSign.GetPubKey().(*kMultiSig.LegacyAminoPubKey)

	willUseSigs := make([]signing.SignatureV2, 0)
	for _, s := range signatures {
		ss, err := c.clientCtx.TxConfig.UnmarshalSignatureJSON(s)
		if err != nil {
			return nil, err
		}
		willUseSigs = append(willUseSigs, ss...)
	}

	multiSigData := multisig.NewMultisig(len(legacyPubkey.PubKeys))
	for _, sig := range willUseSigs {
		if err := multisig.AddSignatureV2(multiSigData, sig, legacyPubkey.GetPubKeys()); err != nil {
			return nil, err
		}
	}

	sigV2 := signing.SignatureV2{
		PubKey:   legacyPubkey,
		Data:     multiSigData,
		Sequence: accountMultiSign.GetSequence(),
	}

	tx, err := c.clientCtx.TxConfig.TxJSONDecoder()(rawTx)
	if err != nil {
		return nil, err
	}
	txBuilder, err := c.clientCtx.TxConfig.WrapTxBuilder(tx)
	if err != nil {
		return nil, err
	}
	txBuilder.SetSignatures(sigV2)
	return c.clientCtx.TxConfig.TxEncoder()(txBuilder.GetTx())
}
