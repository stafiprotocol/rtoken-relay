package config

const (
	LiquidityBondModuleId   = "RTokenSeries"
	LiquidityBondEventId    = "LiquidityBond"
	StorageBondRecords      = "BondRecords"
	MethodExecuteBondRecord = "RTokenSeries.execute_bond_record"

	RtokenVoteModuleId         = "RTokenVotes"
	StorageVotes               = "Votes"
	MethodRacknowledgeProposal = "RTokenVotes.acknowledge_proposal"

	RTokenLedgerModuleId   = "RTokenLedger"
	EraPoolUpdatedEventId  = "EraPoolUpdated"
	StorageChainEras       = "ChainEras"
	MethodSetChainEra      = "RTokenLedger.set_chain_era"
	MethodInitLastVoter    = "RTokenLedger.init_last_voter"
	MethodBondReport       = "RTokenLedger.bond_report"
	MethodActiveReport     = "RTokenLedger.active_report"
	MethodWithdrawReport   = "RTokenLedger.withdraw_report"
	MethodTransferReport   = "RTokenLedger.transfer_report"
	BondReportEventId      = "BondReport"
	WithdrawUnbondEventId  = "WithdrawUnbond"
	TransferBackEventId    = "TransferBack"
	StorageSubAccounts     = "SubAccounts"
	StorageMultiThresholds = "MultiThresholds"
	StorageBondedPools     = "BondedPools"
	StorageSnapshots       = "Snapshots"
	StoragePoolUnbonds     = "PoolUnbonds"
)
