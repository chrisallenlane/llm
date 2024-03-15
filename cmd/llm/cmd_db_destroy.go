package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/chrisallenlane/llm/internal/config"
	"github.com/chrisallenlane/llm/internal/session"
	"gorm.io/gorm"
)

func cmdDBDestroy(opts map[string]interface{}, _ session.Session, _ *gorm.DB) {
	// get the database path
	dbPath, err := config.DBPath(opts, runtime.GOOS)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to construct database path: %s: %v\n", dbPath, err)
		os.Exit(1)
	}

	// delete the database file
	if err := os.Remove(dbPath); err != nil {
		fmt.Fprintf(os.Stderr, "failed to delete database: %s: %v\n", dbPath, err)
		os.Exit(1)
	}
}
