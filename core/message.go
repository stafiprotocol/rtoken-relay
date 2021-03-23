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

	InitLastVoter    = Reason("InitLastVoter")
	NewEra           = Reason("NewEra")
	BondedPools      = Reason("BondedPools")
	EraPoolUpdated   = Reason("EraPoolUpdated")
	BondReport       = Reason("BondReport")
	ActiveReport     = Reason("ActiveReport")
	BondReportEvent  = Reason("BondReportEvent")
	NewMultisig      = Reason("AsMulti")
	MultisigExecuted = Reason("MultisigExecuted")
)
