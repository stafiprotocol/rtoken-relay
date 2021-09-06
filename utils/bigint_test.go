package utils

import (
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/assert"
)

var oneEth = big.NewInt(1000000000000000000)

func TestFromString(t *testing.T) {
	a := "32000000000000000000"
	b, ok := StringToBigint(a)
	assert.Equal(t, true, ok)

	x := big.NewInt(32)
	x.Mul(x, oneEth)
	assert.Equal(t, 0, b.Cmp(x))
}

func TestByteToBigInt(t *testing.T) {
	a := big.NewInt(100)

	t.Log(hexutil.Encode(a.Bytes()))

	b := a.Bytes()
	c := big.NewInt(0).SetBytes(b)

	assert.Equal(t, a.Uint64(), c.Uint64())
}

func TestTimeNow(t *testing.T) {
	currentTime := time.Now()
	t.Log(currentTime.Unix())
	newtime := currentTime.Add(1 * time.Hour)
	t.Log(newtime.Unix())
}

func TestMap(t *testing.T) {
	bz := []byte(`abc`)

	a := map[string][]byte{"x": bz}

	bz = append(bz, 0x01)

	t.Log(a)
	t.Log(bz)
}
