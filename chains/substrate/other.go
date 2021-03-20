package substrate

import (
	"github.com/stafiprotocol/rtoken-relay/config"
	"github.com/stafiprotocol/rtoken-relay/core"
	"time"
)

func (c *Connection) submitSignature(param *core.SubmitSignatureParams) bool {
	for i := 0; i < BlockRetryLimit; i++ {
		c.log.Info("submitSignature on chain...")
		ext, err := c.gc.NewUnsignedExtrinsic(config.SubmitSignatures, param.Symbol, param.Era, param.Pool,
			param.TxType, param.ProposalId, param.Signature)
		err = c.gc.SignAndSubmitTx(ext)
		if err != nil {
			if err.Error() == TerminatedError.Error() {
				c.log.Error("submitSignature  met TerminatedError")
				return false
			}
			c.log.Error("submitSignature error", "err", err)
			time.Sleep(BlockRetryInterval)
			continue
		}
		return true
	}
	return true
}
