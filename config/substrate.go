package config

const (
	TransferModuleId  = "Balances"
	TransferKeepAlive = "transfer_keep_alive"
	Transfer          = "transfer"

	StakingModuleId  = "Staking"
	StorageActiveEra = "ActiveEra"

	StorageLegder        = "Ledger"
	MethodAsMulti        = "Multisig.as_multi"
	MethodApproveAsMulti = "Multisig.approve_as_multi"

	MethodUnbond    = "Staking.unbond"
	MethodBondExtra = "Staking.bond_extra"

	MultisigModuleId        = "Multisig"
	NewMultisigEventId      = "NewMultisig"
	MultisigExecutedEventId = "MultisigExecuted"
	StorageMultisigs        = "Multisigs"

	ParamDest     = "dest"
	ParamDestType = "Address"

	ParamValue     = "value"
	ParamValueType = "Compact<Balance>"
)
