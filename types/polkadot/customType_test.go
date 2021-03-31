package polkadot

import (
	"testing"
)

func TestRegCustomTypes(t *testing.T) {
	RuntimeType{}.Reg()
}
