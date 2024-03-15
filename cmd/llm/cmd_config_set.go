package main

import (
	"fmt"
	"os"

	"github.com/chrisallenlane/llm/internal/config"
	"github.com/chrisallenlane/llm/internal/session"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func cmdConfigSet(opts map[string]interface{}, _ session.Session, db *gorm.DB) {
	conf := config.Config{
		Name:  opts["<key>"].(string),
		Value: opts["<val>"].(string),
	}

	result := db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "name"}},
		DoUpdates: clause.AssignmentColumns([]string{"value"}),
	}).Create(&conf)
	if result.Error != nil {
		fmt.Fprintf(os.Stderr, "failed to set config: %v\n", result.Error)
		os.Exit(1)
	}
}
