package utils

import (
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stafiprotocol/rtoken-relay/models/submodel"
)

func StrReceives(receives []*submodel.Receive) string {
	ret := "["
	for _, r := range receives {
		ret += hexutil.Encode(r.Recipient)
		ret += ", "
	}
	if len(ret) > 3 {
		ret = ret[:len(ret)-2]
	}
	ret += "]"
	return ret
}
