package main

import (
	"fmt"
	"os"

	"github.com/chrisallenlane/llm/internal/config"
	"github.com/chrisallenlane/llm/internal/session"
	"gorm.io/gorm"
)

func cmdConfigList(_ map[string]interface{}, _ session.Session, db *gorm.DB) {
	var configs []config.Config

	// query configs from the database
	if err := db.Find(&configs).Error; err != nil {
		fmt.Fprintf(os.Stderr, "failed to load configs: %v\n", err)
		os.Exit(1)
	}

	// write to stdout
	for _, c := range configs {
		// TODO: lay out in columns
		fmt.Printf("%s\t%s\n", c.Name, c.Value)
	}
}
