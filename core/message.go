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

	CurrentChainEra = Reason("CurrentChainEra")
	NewEra                 = Reason("NewEra")
	BondedPools            = Reason("BondedPools")
	EraPoolUpdated         = Reason("EraPoolUpdated")
	InformChain            = Reason("InformChain")
	BondReportEvent        = Reason("BondReportEvent")
	ActiveReport           = Reason("ActiveReport")
	ActiveReportedEvent    = Reason("ActiveReportedEvent")
	WithdrawReportedEvent  = Reason("WithdrawReportedEvent")
	TransferReportedEvent  = Reason("TransferReportedEvent")
	NominationUpdatedEvent = Reason("NominationUpdatedEvent")
	NewMultisig            = Reason("AsMulti")
	MultisigExecuted       = Reason("MultisigExecuted")
	GetEraNominated        = Reason("GetEraNominated")
	//cosmos use
	SubmitSignature       = Reason("SubmitSignature")
	SignatureEnough       = Reason("SignatureEnough")
	ValidatorUpdatedEvent = Reason("ValidatorUpdatedEvent")
)
