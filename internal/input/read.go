package input

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

// Read returns the user input, regardless of how it was passed
func Read(opt string, opts map[string]interface{}) (string, error) {
	prompt := ""

	// get stdin mode
	stat, err := os.Stdin.Stat()
	if err != nil {
		return "", fmt.Errorf("failed to open stdin: %v", err)
	}

	// if the prompt was passed via an opt, read it
	if opts[opt] != nil {
		prompt = opts[opt].(string)

		// read input on `stdin` if available
	} else if (stat.Mode() & os.ModeCharDevice) == 0 {
		reader := bufio.NewReader(os.Stdin)

		for {
			line, err := reader.ReadString('\n')
			if err != nil && err == io.EOF {
				break
			} else if err != nil {
				return "", fmt.Errorf("failed to read stdin: %v", err)
			}
			prompt += line
		}

		// otherwise, open the user's editor
	} else {
		var err error
		prompt, err = Editor("", opts)
		if err != nil {
			return "", fmt.Errorf("failed to read from user editor: %v", err)
		}
	}

	return strings.TrimSpace(prompt), nil
}
