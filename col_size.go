package fancy_printer

import (
	"os"

	"golang.org/x/crypto/ssh/terminal"
)

func getTerminalSize() (int, error) {
	terminalWidth, _, err := terminal.GetSize(int(os.Stdout.Fd()))
	return terminalWidth, err
}

func getTerminalSizeMock() (int, error) {
	terminalWidth := 120
	return terminalWidth, nil
}
