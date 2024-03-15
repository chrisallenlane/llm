package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/chrisallenlane/llm/internal/session"
	"gorm.io/gorm"
)

// view the session information
func cmdSessionView(opts map[string]interface{}, _ session.Session, db *gorm.DB) {
	// get the session name
	name := opts["<name>"].(string)

	// read the session from the database
	var s session.Session
	result := db.Where("name = ?", name).Find(&s)
	if result.Error != nil {
		fmt.Fprintf(os.Stderr, "failed to find session %s: %v\n", name, result.Error)
		os.Exit(1)
	}
	if result.RowsAffected == 0 {
		fmt.Fprintf(os.Stderr, "session does not exist: %s\n", name)
		os.Exit(1)
	}

	// print the session hint to stdout
	fmt.Println(strings.Trim(s.Hint, "\n"))
}
