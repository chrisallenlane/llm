package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/chrisallenlane/llm/internal/config"
	"github.com/chrisallenlane/llm/internal/session"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func cmdSessionUse(opts map[string]interface{}, _ session.Session, db *gorm.DB) {
	// get the session name
	name := strings.TrimSpace(opts["<name>"].(string))
	if name == "" {
		fmt.Fprintf(os.Stderr, "<name> must not be empty.")
		os.Exit(1)
	}

	// verify that the session exists
	result := db.Where("name = ?", name).First(&session.Session{})
	if result.Error != nil {
		fmt.Fprintf(os.Stderr, "failed to load session %s: %v\n", name, result.Error)
		os.Exit(1)
	}
	if result.RowsAffected == 0 {
		fmt.Fprintf(os.Stderr, "session %s does not exist.\n", name)
		os.Exit(1)
	}

	// use the new session
	conf := config.Config{
		Name:  "session",
		Value: name,
	}

	result = db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "name"}},
		DoUpdates: clause.AssignmentColumns([]string{"value"}),
	}).Create(&conf)
	if result.Error != nil {
		fmt.Fprintf(os.Stderr, "failed to update session: %v\n", result.Error)
		os.Exit(1)
	}
}
