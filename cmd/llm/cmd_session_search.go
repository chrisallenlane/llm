package main

import (
	"fmt"
	"os"

	"github.com/chrisallenlane/llm/internal/session"
	"gorm.io/gorm"
)

func cmdSessionSearch(opts map[string]interface{}, sess session.Session, db *gorm.DB) {
	// construct the user-specified LIKE qualifier
	like := "%" + opts["<string>"].(string) + "%"

	// query out matching messages
	err := db.Model(&sess).
		Where("content LIKE ?", like).
		Association("Messages").
		Find(&sess.Messages)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load messages: %v\n", err)
		os.Exit(1)
	}

	// display the matching messages
	for _, msg := range sess.Messages {
		msg.Display()
	}
}
