package core

type Message struct {
	Source      RSymbol
	Destination RSymbol
	Reason      Reason
	Content     interface{}
}

type Reason string

const (
	LiquidityBond       = Reason("LiquidityBond")
	LiquidityBondResult = Reason("LiquidityBondResult")

	NewEra           = Reason("NewEra")
	EraPoolUpdated   = Reason("EraPoolUpdated")
	NewMultisig      = Reason("AsMulti")
	MultisigExecuted = Reason("MultisigExecuted")
)
