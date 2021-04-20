package config

const (
	RTokenSeriesModuleId     = "RTokenSeries"
	LiquidityBondEventId     = "LiquidityBond"
	NominationUpdatedEventId = "NominationUpdated"
	StorageBondRecords       = "BondRecords"
	MethodExecuteBondRecord  = "RTokenSeries.execute_bond_record"

	RtokenVoteModuleId         = "RTokenVotes"
	StorageVotes               = "Votes"
	MethodRacknowledgeProposal = "RTokenVotes.acknowledge_proposal"

	RTokenLedgerModuleId       = "RTokenLedger"
	EraPoolUpdatedEventId      = "EraPoolUpdated"
	StorageChainEras           = "ChainEras"
	StorageCurrentEraSnapShots = "CurrentEraSnapShots"
	MethodSetChainEra          = "RTokenLedger.set_chain_era"
	MethodBondReport           = "RTokenLedger.bond_report"
	MethodActiveReport         = "RTokenLedger.active_report"
	MethodWithdrawReport       = "RTokenLedger.withdraw_report"
	MethodTransferReport       = "RTokenLedger.transfer_report"
	BondReportedEventId        = "BondReported"
	ActiveReportedEventId      = "ActiveReported"
	WithdrawReportedEventId    = "WithdrawReported"
	StorageSubAccounts         = "SubAccounts"
	StorageMultiThresholds     = "MultiThresholds"
	StorageBondedPools         = "BondedPools"
	StorageSnapshots           = "Snapshots"
	StoragePoolUnbonds         = "PoolUnbonds"
	SignaturesEnoughEventId    = "SignaturesEnough"
	StorageSignatures          = "Signatures"
)
