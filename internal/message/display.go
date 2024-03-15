package message

import (
	"fmt"
)

// Display prints a message to stdout
func (msg *Message) Display() {
	// create a subheading
	header := fmt.Sprintf(
		"## id: %d | role: %s | date: %s",
		msg.ID,
		msg.Role,
		msg.Date,
	)
	// display the log output
	fmt.Println(Format(header))
	fmt.Println(Format(msg.Content))
}
