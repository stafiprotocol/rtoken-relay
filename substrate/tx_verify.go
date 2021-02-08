package substrate

import (
	"github.com/stafiprotocol/rtoken-relay/conn"
)

func (sc *SarpcClient) TransferVerify(record *conn.BondRecord) (bool, conn.OpposeReason) {
	return true, conn.NoReason
}
