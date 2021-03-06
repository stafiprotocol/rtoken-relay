package config

const (
	LiquidityBondModuleId = "RTokenSeries"
	LiquidityBondEventId  = "LiquidityBond"
	StorageBondRecords    = "BondRecords"
	ExecuteBondRecord     = "RTokenSeries.execute_bond_record"
	RtokenVoteModuleId    = "RTokenVotes"
	StorageVotes          = "Votes"
	RacknowledgeProposal  = "RTokenVotes.acknowledge_proposal"
	RTokenLedgerModuleId  = "RTokenLedger"
	EraUpdatedEventId     = "EraPoolUpdated"
	StorageChainEras      = "ChainEras"
	SetChainEra           = "RTokenLedger.set_chain_era"
	SetPoolActive         = "RTokenLedger.set_pool_active"
	//PoolSubAccountAddedEventId = "PoolSubAccountAdded"
	//StoragePoolBonded          = "PoolBonded"
	//StoragePoolSubAccountFlag  = "PoolSubAccountFlag"
	StorageTotalLinking = "TotalLinking"
	StorageBondFaucets  = "BondFaucets"
)