package config

const (
	BalancesModuleId        = "Balances"
	TransferKeepAlive       = "transfer_keep_alives"
	Transfer                = "transfer"
	MethodTransferKeepAlive = "Balances.transfer_keep_alive"
	ConstExistentialDeposit = "ExistentialDeposit"

	StakingModuleId         = "Staking"
	StorageActiveEra        = "ActiveEra"
	StorageNominators       = "Nominators"
	StorageErasRewardPoints = "ErasRewardPoints"
	StorageBonded           = "Bonded"
	StorageLedger           = "Ledger"
	MethodPayoutStakers     = "Staking.payout_stakers"
	MethodUnbond            = "Staking.unbond"
	MethodBondExtra         = "Staking.bond_extra"
	MethodWithdrawUnbonded  = "Staking.withdraw_unbonded"

	MultisigModuleId        = "Multisig"
	NewMultisigEventId      = "NewMultisig"
	MultisigExecutedEventId = "MultisigExecuted"
	StorageMultisigs        = "Multisigs"
	MethodAsMulti           = "Multisig.as_multi"

	SystemModuleId = "System"
	StorageAccount = "Account"

	MethodBatch = "Utility.batch"

	ParamDest     = "dest"
	ParamDestType = "Address"

	ParamValue     = "value"
	ParamValueType = "Compact<Balance>"
)
