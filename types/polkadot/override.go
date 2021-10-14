package polkadot

import "strings"

var typesModules = map[string]map[string]string{
	"parasInclusion": {
		"validatorindex": "ParaValidatorIndex",
	},
	"inclusion": {
		"validatorindex": "ParaValidatorIndex",
	},
	"parasScheduler": {
		"validatorindex": "ParaValidatorIndex",
	},
	"parasShared": {
		"validatorindex": "ParaValidatorIndex",
	},
	"scheduler": {
		"validatorindex": "ParaValidatorIndex",
	},
	"paraDisputes": {
		"validatorindex": "ParaValidatorIndex",
	},
	"parasDisputes": {
		"validatorindex": "ParaValidatorIndex",
	},
	"shared": {
		"validatorindex": "ParaValidatorIndex",
	},
	"assets": {
		"approval":       "AssetApproval",
		"approvalkey":    "AssetApprovalKey",
		"balance":        "TAssetBalance",
		"destroywitness": "AssetDestroyWitness",
	},
	"beefy": {
		"authorityid": "BeefyId",
	},
	"ethereum": {
		"block":             "EthBlock",
		"header":            "EthHeader",
		"receipt":           "EthReceipt",
		"transaction":       "EthTransaction",
		"transactionstatus": "EthTransactionStatus",
	},
}

func (r *RuntimeType) overrideModuleType(t string) string {
	moduleTypes, ok := typesModules[strings.ToLower(r.Module)]
	if !ok {
		return t
	}
	if overrideType, ok := moduleTypes[t]; ok {
		return overrideType
	}
	return t
}
