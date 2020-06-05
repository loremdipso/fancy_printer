package fancy_printer

import (
	"fmt"
	"os"

	"github.com/loremdipso/go_utils"

	// TODO: this maaaaay be overkill. Just need the terminal size
	"github.com/fatih/color"
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

	// TODO: we need to duplicate the array since it needs to be sorted.
	// Make this less fragile
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

func getTagCols(tokens []string) ([][]string, int, error) {
	sortedTags := go_utils.SortByLength(go_utils.DupStrArray(tokens))
	numCols, err := getNumTagCols(sortedTags)
	if err != nil {
		return nil, 0, err
	}
	// numCols -= 1 // TODO: this is generally wrong. Fix it
	// numCols -= 1

	cols := make([][]string, numCols)
	numRows := len(sortedTags)/numCols + 1
	leftOvers := len(sortedTags) - numCols*numRows

	count := 0
	for i := 0; i < numCols; i++ {
		cols[i] = make([]string, numRows+1)
		localNumRows := numRows
		if i < leftOvers {
			localNumRows += 1
		}

		for j := 0; j <= localNumRows; j++ {
			tagIndex := i*localNumRows + j
			if tagIndex >= len(sortedTags) {
				break
			}
			cols[i][j] = sortedTags[tagIndex]
			count++
		}
	}

	if leftOvers > 0 {
		numRows += 1
	}
	return cols, numRows, nil
}

func getNumTagCols(sortedTags []string) (int, error) {
	terminalWidth, _, err := terminal.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return 0, err
	}

	numCols := 1
	for {
		// test if valid. If not, return old
		total := 0
		numCols++
		rowsPerCol := len(sortedTags)/numCols + 1
		if rowsPerCol < 1 {
			return numCols - 1, nil
		}

		for col := 0; col < numCols; col++ {
			// want the largest-lengthed string in column, which is the last element
			tagIndexStart := go_utils.Min(len(sortedTags), (col * rowsPerCol))
			tagIndexEnd := go_utils.Min((tagIndexStart + rowsPerCol - 1), len(sortedTags))
			longestInCol := go_utils.FindLongest(sortedTags[tagIndexStart:tagIndexEnd])
			total += longestInCol
		}
		total += (numCols - 1) * len(columnSeparator)
		if total > terminalWidth {
			return numCols - 1, nil
		}
	}
}
