package polkadot

import (
	"github.com/itering/scale.go/types"
	"testing"
)

func TestRegCustomTypes(t *testing.T) {
	types.RuntimeType{}.Reg()
}
