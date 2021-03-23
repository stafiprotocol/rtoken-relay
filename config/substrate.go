package config

const (
	TransferModuleId  = "Balances"
	TransferKeepAlive = "transfer_keep_alive"
	Transfer          = "transfer"

	StakingModuleId         = "Staking"
	StorageActiveEra        = "ActiveEra"
	StorageNominators       = "Nominators"
	StorageErasRewardPoints = "ErasRewardPoints"
	StorageBonded           = "Bonded"
	StorageLedger           = "Ledger"
	MethodPayoutStakers     = "Staking.payout_stakers"
	MethodUnbond            = "Staking.unbond"
	MethodBondExtra         = "Staking.bond_extra"

	MethodAsMulti        = "Multisig.as_multi"
	MethodApproveAsMulti = "Multisig.approve_as_multi"

	MultisigModuleId        = "Multisig"
	NewMultisigEventId      = "NewMultisig"
	MultisigExecutedEventId = "MultisigExecuted"
	StorageMultisigs        = "Multisigs"

	MethodBatch = "Utility.batch"

	ParamDest     = "dest"
	ParamDestType = "Address"

	ParamValue     = "value"
	ParamValueType = "Compact<Balance>"
)
