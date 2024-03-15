package main

import (
	"fmt"
	"os"

	"github.com/chrisallenlane/llm/internal/session"
	"gorm.io/gorm"
)

func cmdSessionRemove(opts map[string]interface{}, _ session.Session, db *gorm.DB) {
	name := opts["<name>"].(string)

	result := db.Where("name = ?", name).Delete(&session.Session{})

	if result.Error != nil {
		fmt.Fprintf(os.Stderr, "cannot remove session %s: %v\n", name, result.Error)
		os.Exit(1)
	}
}
