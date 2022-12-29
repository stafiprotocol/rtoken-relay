package utils

import (
	"fmt"
	"testing"

	bncCmnTypes "github.com/stafiprotocol/go-sdk/common/types"
	"github.com/stretchr/testify/assert"
)

var (
	csvFile = "../datas/test.csv"
)

func TestWriteAndReadeCSV(t *testing.T) {
	contents := [][]string{
		{"id1", "name1", "score60"},
		{"id2", "name2", "score62"},
	}

	err := WriteCSV(csvFile, contents)
	assert.NoError(t, err)

	lines := ReadCSV(csvFile)
	t.Log(lines)
}

func TestWriteAndReadeCSV1(t *testing.T) {
	contents := [][]string{
		{"id1", "name1", "score63"},
		{"id2", "name2", "score64"},
	}

	err := WriteCSV(csvFile, contents)
	assert.NoError(t, err)

	lines := ReadCSV(csvFile)
	t.Log(lines)
}

func TestGetStakeCAoB(t *testing.T) {
	exp, err := bncCmnTypes.AccAddressFromHex("91D7deA99716Cbb247E81F1cfB692009164a967E")
	if err != nil {
		t.Fatal(err)
	}
	stakeCAoB := GetStakeCAoB(exp.Bytes(), DelegateCAoBSalt)
	fmt.Println(stakeCAoB.String())
	if delAddr := GetStakeCAoB(stakeCAoB.Bytes(), DelegateCAoBSalt); delAddr.String() != exp.String() {
		t.Fatal()
	}
}
