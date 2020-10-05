package fancy_printer

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/loremdipso/go_utils"

	// TODO: this maaaaay be overkill.
	// we really just need the terminal size
	"golang.org/x/crypto/ssh/terminal"
)

var columnSeparator = " | "
var columnColors = []func(string, ...interface{}) string{
	color.HiRedString,
	color.HiBlueString,
	color.HiGreenString,
	color.HiYellowString,
	color.HiWhiteString,
}

func GetTruncatedLine(prefix string, lineToPrint string) (string, string, error) {
	terminalWidth, _, err := terminal.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return "", "", err
	}

	remainingWidth := terminalWidth - len(prefix)
	ellipsis := ""
	if len(lineToPrint) > remainingWidth {
		ellipsis = "..."
		remainingWidth -= 3
	}

	return lineToPrint[:go_utils.Min(len(lineToPrint), remainingWidth)], ellipsis, nil
}

func PrintArrayAsGrid(incomingTokens []string, simple bool, useColors bool) error {
	if len(incomingTokens) == 0 {
		return nil
	}

	// We need to duplicate the array since it needs to be sorted.
	// TODO: make this less fragile
	tokens := go_utils.DupStrArray(incomingTokens)

	if simple {
		fmt.Println(go_utils.StringArrayToString(tokens))
	} else {
		// The idea is that we want to minimize the amount of space the tokens take up
		// we'll do this by sorting them by size, popping them off each end as we go along

		cols, numRows, err := getTagCols(tokens)
		if err != nil {
			return err
		}

		for rowIndex := 0; rowIndex < numRows; rowIndex++ {
			// TODO: use string builder? Or is it fastest to print to output?
			// for now anything is fine. Performance isn't a terribly important consideration.
			var line string

			for columnIndex, column := range cols {
				tag := column[rowIndex]
				columnWidth := go_utils.FindLongest(column)
				if len(tag) > 0 {
					if useColors {
						color := columnColors[go_utils.Mod(columnIndex, len(columnColors))]
						line += fmt.Sprint(color("%-*s", columnWidth, tag))
					} else {
						line += fmt.Sprintf("%-*s", columnWidth, tag)
					}

					// only add separator if there's something in the next row
					if columnIndex < len(cols)-1 && cols[columnIndex+1][rowIndex] != "" {
						line += columnSeparator
					}
				}
			}

			fmt.Println(line)
		}
	}
	return nil
}

// returns the input data as an array of columns, the maximum number of columns, or an error
// if something went squiffy
func getTagCols(tokens []string) ([][]string, int, error) {
	terminalWidth, err := getTerminalSize()

	if err != nil {
		return nil, 0, err
	}

	sortedTokens := go_utils.SortByLength(go_utils.DupStrArray(tokens))
	oldCols, _ := splitIntoNColumns(sortedTokens, 1)
	numCols := 2
	for {
		newCols, haveEmptyColumns := splitIntoNColumns(sortedTokens, numCols)
		if haveEmptyColumns {
			break
		}

		lengthOfLongestRow := 0
		for _, column := range newCols {
			longestInCol := go_utils.FindLongest(column)
			lengthOfLongestRow += longestInCol
		}

		lengthOfLongestRow += (numCols - 1) * len(columnSeparator)
		if lengthOfLongestRow > terminalWidth {
			break
		} else {
			oldCols = newCols
			numCols++
		}
	}

	numRows := 0
	if len(oldCols) > 0 {
		numRows = len(oldCols[0])
	}
	return oldCols, numRows, nil
}

func splitIntoNColumns(sortedTokens []string, numCols int) ([][]string, bool) {
	numRows := len(sortedTokens) / numCols
	hasLeftovers := (numRows * numCols) < len(sortedTokens)
	hasEmptyCols := false
	if hasLeftovers {
		numRows += 1
	}

	cols := make([][]string, numCols)
	index := 0
	for c := 0; c < numCols; c++ {
		cols[c] = make([]string, numRows)
		for r := 0; r < numRows; r++ {
			if index >= len(sortedTokens) {
				if c < numCols-1 || r == 0 {
					hasEmptyCols = true
				}
				return cols, hasEmptyCols
			}

			cols[c][r] = sortedTokens[index]
			index++
		}
	}

	return cols, hasEmptyCols
}
