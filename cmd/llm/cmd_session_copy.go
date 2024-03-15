package main

import (
	"fmt"
	"os"

	"github.com/chrisallenlane/llm/internal/message"
	"github.com/chrisallenlane/llm/internal/session"
	"gorm.io/gorm"
)

func cmdSessionCopy(opts map[string]interface{}, _ session.Session, db *gorm.DB) {
	// get the <orig> session
	var orig session.Session
	result := db.Preload("Messages").
		Where("name = ?", opts["<orig>"].(string)).
		Limit(1).
		Find(&orig)
	if result.Error != nil {
		fmt.Fprintf(os.Stderr, "failed to load session: %v\n", result.Error)
		os.Exit(1)
	}

	if result.RowsAffected == 0 {
		fmt.Fprintf(os.Stderr, "could not find session: %v\n", result.Error)
		os.Exit(1)
	}

	// copy the session
	name := opts["<name>"].(string)
	dupe := orig
	dupe.ID = 0
	dupe.Name = name

	// if the user provided the `--all` flag, copy the messages too
	if opts["--all"] != nil {
		dupe.Messages = make([]message.Message, len(orig.Messages))
		for i, msg := range orig.Messages {
			dupe.Messages[i] = msg
			dupe.Messages[i].ID = 0
		}
	}

	// save the session
	result = db.Save(&dupe)
	if result.Error != nil {
		fmt.Fprintf(os.Stderr, "failed to save session %s: %v\n", name, result.Error)
		os.Exit(1)
	}
}
