package conn

type Validator interface {
	TransferVerify(record *BondRecord) (bool, OpposeReason)
}
