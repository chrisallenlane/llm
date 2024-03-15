package main

import (
	"fmt"
	"os"

	"github.com/chrisallenlane/llm/internal/input"
	"github.com/chrisallenlane/llm/internal/session"
	"gorm.io/gorm"
)

func cmdSessionEdit(opts map[string]interface{}, sess session.Session, db *gorm.DB) {
	// get the session name
	name := opts["<name>"].(string)

	// read the session from the database
	result := db.Where("name = ?", name).Find(&sess)
	if result.Error != nil {
		fmt.Fprintf(os.Stderr, "failed to find session %s: %v\n", name, result.Error)
		os.Exit(1)
	}
	if result.RowsAffected == 0 {
		fmt.Fprintf(os.Stderr, "session does not exist: %s", name)
		os.Exit(1)
	}

	// update the session hint
	hint, err := input.Editor(sess.Hint, opts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open editor: %v\n", err)
		os.Exit(1)
	}
	sess.Hint = hint

	result = db.Save(&sess)
	if result.Error != nil {
		fmt.Fprintf(os.Stderr, "failed to save session %s: %v\n", name, result.Error)
		os.Exit(1)
	}
}
