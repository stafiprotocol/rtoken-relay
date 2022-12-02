package core

type Message struct {
	Source      RSymbol
	Destination RSymbol
	Reason      Reason
	Content     interface{}
}

type Reason string

const (

	// send by other chain
	LiquidityBondResult        = Reason("LiquidityBondResult")
	ExeLiquidityBondAndSwap    = Reason("ExeLiquidityBondAndSwap")
	NewEra                     = Reason("NewEra")
	InformChain                = Reason("InformChain")
	ActiveReport               = Reason("ActiveReport")
	GetEraNominated            = Reason("GetEraNominated")
	GetBondState               = Reason("GetBondState")
	WaitAndGetSubmitSignatures = Reason("WaitAndGetSubmitSignatures")

	// send by stafi
	BondedPools = Reason("BondedPools")

	// stafi event, send by stafi
	LiquidityBondEvent     = Reason("LiquidityBondEvent")
	EraPoolUpdatedEvent    = Reason("EraPoolUpdatedEvent")
	BondReportedEvent      = Reason("BondReportedEvent")
	ActiveReportedEvent    = Reason("ActiveReportedEvent")
	WithdrawReportedEvent  = Reason("WithdrawReportedEvent")
	TransferReportedEvent  = Reason("TransferReportedEvent")
	NominationUpdatedEvent = Reason("NominationUpdatedEvent")

	// other substrate chain multisig event
	NewMultisig      = Reason("AsMulti")
	MultisigExecuted = Reason("MultisigExecuted")

	// matic use
	SubmitSignature      = Reason("SubmitSignature")
	SignatureEnoughEvent = Reason("SignatureEnoughEvent")

	ValidatorUpdatedEvent = Reason("ValidatorUpdatedEvent")
)
