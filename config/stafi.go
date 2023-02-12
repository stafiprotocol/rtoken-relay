package config

const (
	RTokenSeriesModuleId     = "RTokenSeries"
	LiquidityBondEventId     = "LiquidityBond"
	NominationUpdatedEventId = "NominationUpdated"
	ValidatorUpdatedEventId  = "ValidatorUpdated"
	StorageBondRecords       = "BondRecords"
	StorageBondStates        = "BondStates"
	MethodExecuteBondRecord  = "RTokenSeries.execute_bond_record"
	MethodExecuteBondAndSwap = "RTokenSeries.execute_bond_and_swap"
	StorageNominated         = "Nominated"

	RtokenVoteModuleId         = "RTokenVotes"
	StorageVotes               = "Votes"
	MethodRacknowledgeProposal = "RTokenVotes.acknowledge_proposal"

	RTokenLedgerModuleId                      = "RTokenLedger"
	RTokenRelayersModuleId                    = "Relayers"
	EraPoolUpdatedEventId                     = "EraPoolUpdated"
	StorageChainEras                          = "ChainEras"
	StorageCurrentEraSnapShots                = "CurrentEraSnapShots"
	StorageRelayerThreshold                   = "RelayerThreshold"
	StorageActiveChangeRateLimit              = "ActiveChangeRateLimit"
	StorageEraSnapShots                       = "EraSnapShots"
	StorageLeastBond                          = "LeastBond"
	StoragePendingStake                       = "PendingStake"
	StoragePendingReward                      = "PendingReward"
	MethodSetChainEra                         = "RTokenLedger.set_chain_era"
	MethodBondReport                          = "RTokenLedger.bond_report"
	MethodNewBondReport                       = "RTokenLedger.new_bond_report"
	MethodActiveReport                        = "RTokenLedger.active_report"
	MethodNewActiveReport                     = "RTokenLedger.new_active_report"
	MethodBondAndReportActive                 = "RTokenLedger.bond_and_report_active"
	MethodBondAndReportActiveWithPendingValue = "RTokenLedger.bond_and_report_active_with_pending_value"
	MethodWithdrawReport                      = "RTokenLedger.withdraw_report"
	MethodTransferReport                      = "RTokenLedger.transfer_report"
	BondReportedEventId                       = "BondReported"
	ActiveReportedEventId                     = "ActiveReported"
	WithdrawReportedEventId                   = "WithdrawReported"
	TransferReportedEventId                   = "TransferReported"
	StorageSubAccounts                        = "SubAccounts"
	StorageMultiThresholds                    = "MultiThresholds"
	StorageBondedPools                        = "BondedPools"
	StorageSnapshots                          = "Snapshots"
	StoragePoolUnbonds                        = "PoolUnbonds"
	SignaturesEnoughEventId                   = "SignaturesEnough"
	StorageSignatures                         = "Signatures"
	SubmitSignatures                          = "RTokenSeries.submit_signatures"
)
