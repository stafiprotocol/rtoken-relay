package substrate

import (
	"github.com/stafiprotocol/rtoken-relay/utils"
)

type TransferBack struct {
	Symbol string
	Pool   string
	Era    string
}

func CreateTransferback(filePath string, tb *TransferBack) error {
	lines := utils.ReadCSV(filePath)
	newline := []string{tb.Symbol, tb.Pool, tb.Era}
	lines = append(lines, newline)

	return utils.WriteCSV(filePath, lines)
}

func IsTransferbackExist(filePath string, tb *TransferBack) bool {
	lines := utils.ReadCSV(filePath)
	for _, line := range lines {
		if len(line) != 3 {
			panic("size of line is not 3")
		}
		if line[0] == tb.Symbol && line[1] == tb.Pool && line[2] == tb.Era {
			return true
		}
	}

	return false
}

func DeleteTransferback(filePath string, tb *TransferBack) error {
	lines := utils.ReadCSV(filePath)
	newLines := make([][]string, 0)
	for _, line := range lines {
		if len(line) != 3 {
			panic("size of line is not 3")
		}

		if line[0] != tb.Symbol || line[1] != tb.Pool && line[2] != tb.Era {
			newLines = append(newLines, line)
		}
	}

	return utils.WriteCSV(filePath, newLines)
}
