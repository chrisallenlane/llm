package main

import (
	"fmt"
	"os"

	"gorm.io/gorm"

	"github.com/chrisallenlane/llm/internal/session"
)

func cmdSessionList(_ map[string]interface{}, sess session.Session, db *gorm.DB) {
	var sessions []session.Session

	// query sessions from the database
	if err := db.Find(&sessions).Error; err != nil {
		fmt.Fprintf(os.Stderr, "failed to load sessions: %v\n", err)
		os.Exit(1)
	}

	// write them to stdout
	for _, s := range sessions {
		indicator := "  "
		if s.Name == sess.Name {
			indicator = "* "
		}
		fmt.Println(indicator + s.Name)
	}
}
