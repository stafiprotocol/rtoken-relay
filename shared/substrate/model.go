package substrate

import (
	"errors"

	"github.com/shopspring/decimal"
)

const (
	ChainTypeStafi    = "stafi"
	ChainTypePolkadot = "polkadot"

	AddressTypeAccountId    = "AccountId"
	AddressTypeMultiAddress = "MultiAddress"
)

type ChainExtrinsic struct {
	ID                 uint            `gorm:"primary_key"`
	ExtrinsicIndex     string          `json:"extrinsic_index" sql:"default: null;size:100"`
	BlockNum           int             `json:"block_num" `
	BlockTimestamp     int             `json:"block_timestamp"`
	ExtrinsicLength    string          `json:"extrinsic_length"`
	VersionInfo        string          `json:"version_info"`
	CallCode           string          `json:"call_code"`
	CallModuleFunction string          `json:"call_module_function"  sql:"size:100"`
	CallModule         string          `json:"call_module"  sql:"size:100"`
	Params             interface{}     `json:"params" sql:"type:MEDIUMTEXT;" `
	AccountId          string          `json:"account_id"`
	Signature          string          `json:"signature"`
	Nonce              int             `json:"nonce"`
	Era                string          `json:"era"`
	ExtrinsicHash      string          `json:"extrinsic_hash" sql:"default: null" `
	IsSigned           bool            `json:"is_signed"`
	Success            bool            `json:"success"`
	Fee                decimal.Decimal `json:"fee" sql:"type:decimal(30,0);"`
	BatchIndex         int             `json:"-" gorm:"-"`
}

var (
	TerminatedError           = errors.New("terminated")
	BondEqualToUnbondError    = errors.New("BondEqualToUnbondError")
	DiffSmallerThanLeastError = errors.New("DiffSmallerThanLeastError")
)
