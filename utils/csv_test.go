package utils

import (
	"testing"

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
