package main

import (
	"fmt"
	"os"

	"gorm.io/gorm"

	"github.com/chrisallenlane/llm/internal/input"
	"github.com/chrisallenlane/llm/internal/session"
)

func cmdSessionNew(opts map[string]interface{}, sess session.Session, db *gorm.DB) {
	// open the user's EDITOR and read the hint text
	hint, err := input.Editor("", opts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open editor: %v\n", err)
		os.Exit(1)
	}

	// initialize a session
	sess = session.Session{
		Name: opts["<name>"].(string),
		Hint: hint,
	}

	// save to the database
	if err := db.Create(&sess).Error; err != nil {
		fmt.Fprintf(os.Stderr, "failed to create session: %v\n", err)
		os.Exit(1)
	}
}
