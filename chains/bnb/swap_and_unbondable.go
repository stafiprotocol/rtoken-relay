package bnb

import (
	"github.com/stafiprotocol/rtoken-relay/utils"
)

const (
	FromBsc = string("FromBsc")
	FromBc  = string("FromBc")
)

type Swap struct {
	Symbol string
	Pool   string
	Era    string
	From   string
}

type Unbondable struct {
	Symbol string
	Pool   string
	Era    string
}

func CreateSwap(filePath string, swap *Swap) error {
	lines := utils.ReadCSV(filePath)
	newline := []string{swap.Symbol, swap.Pool, swap.Era, swap.From}
	lines = append(lines, newline)

	return utils.WriteCSV(filePath, lines)
}

func IsSwapExist(filePath string, swap *Swap) bool {
	lines := utils.ReadCSV(filePath)
	for _, line := range lines {
		if len(line) != 4 {
			panic("IsSwapExist size of line is not 4")
		}
		if line[0] == swap.Symbol && line[1] == swap.Pool && line[2] == swap.Era && line[3] == swap.From {
			return true
		}
	}

	return false
}

func DeleteSwap(filePath string, swap *Swap) error {
	lines := utils.ReadCSV(filePath)
	newLines := make([][]string, 0)
	for _, line := range lines {
		if len(line) != 4 {
			panic("DeleteSwap size of line is not 4")
		}

		if line[0] != swap.Symbol || line[1] != swap.Pool || line[2] != swap.Era || line[3] != swap.From {
			newLines = append(newLines, line)
		}
	}

	return utils.WriteCSV(filePath, newLines)
}

func CreateUnbondable(filePath string, u *Unbondable) error {
	lines := utils.ReadCSV(filePath)
	newline := []string{u.Symbol, u.Pool, u.Era}
	lines = append(lines, newline)

	return utils.WriteCSV(filePath, lines)
}

func IsUnbondableExist(filePath string, u *Unbondable) bool {
	lines := utils.ReadCSV(filePath)
	for _, line := range lines {
		if len(line) != 3 {
			panic("IsSwapExist size of line is not 3")
		}
		if line[0] == u.Symbol && line[1] == u.Pool && line[2] == u.Era {
			return true
		}
	}

	return false
}
