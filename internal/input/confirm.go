// Package input encapsulates user input methods and data
package input

import (
	"bufio"
	"os"
	"strings"
)

// Confirm returns true if the user responds in the affirmative
func Confirm() bool {
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.ToLower(strings.TrimSpace(input))

	return input == "y" || input == "yes"
}
