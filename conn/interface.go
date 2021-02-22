package conn

type Validator interface {
	TransferVerify(record *BondRecord) (BondReason, error)
}
