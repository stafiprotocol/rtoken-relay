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

	InitLastVoter          = Reason("InitLastVoter")
	NewEra                 = Reason("NewEra")
	BondedPools            = Reason("BondedPools")
	EraPoolUpdated         = Reason("EraPoolUpdated")
	InformChain            = Reason("InformChain")
	BondReportEvent        = Reason("BondReportEvent")
	ActiveReport           = Reason("ActiveReport")
	ActiveReportedEvent    = Reason("ActiveReportedEvent")
	WithdrawReportedEvent  = Reason("WithdrawReportedEvent")
	NominationUpdatedEvent = Reason("NominationUpdatedEvent")
	NewMultisig            = Reason("AsMulti")
	MultisigExecuted       = Reason("MultisigExecuted")
	GetEraNominated        = Reason("GetEraNominated")
)
