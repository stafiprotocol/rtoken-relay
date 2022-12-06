package config

const (
	BalancesModuleId        = "Balances"
	TransferKeepAlive       = "transfer_keep_alive"
	Transfer                = "transfer"
	MethodTransferKeepAlive = "Balances.transfer_keep_alive"
	MethodTransfer          = "Balances.transfer"
	ConstExistentialDeposit = "ExistentialDeposit"

	StakingModuleId           = "Staking"
	StorageActiveEra          = "ActiveEra"
	StorageNominators         = "Nominators"
	StorageErasRewardPoints   = "ErasRewardPoints"
	StorageErasStakersClipped = "ErasStakersClipped"
	StorageEraNominated       = "EraNominated"
	StorageBonded             = "Bonded"
	StorageLedger             = "Ledger"
	MethodPayoutStakers       = "Staking.payout_stakers"
	MethodUnbond              = "Staking.unbond"
	MethodBondExtra           = "Staking.bond_extra"
	MethodWithdrawUnbonded    = "Staking.withdraw_unbonded"
	MethodNominate            = "Staking.nominate"

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
