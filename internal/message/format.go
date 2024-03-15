package message

import (
	"os"

	"github.com/charmbracelet/glamour"
	"github.com/mattn/go-isatty"
)

// Format attempts to word-wrap at 80 columns and render markdown
func Format(s string) string {
	if !isatty.IsTerminal(os.Stdout.Fd()) {
		return s
	}

	rendered, err := glamour.Render(s, "dark")
	if err != nil {
		return s
	}

	return rendered
}
