package main

import (
	"fmt"
	"os"

	"github.com/chrisallenlane/llm/internal/message"
	"github.com/chrisallenlane/llm/internal/session"
	"gorm.io/gorm"
)

func cmdQuestionRemove(opts map[string]interface{}, _ session.Session, db *gorm.DB) {
	name := opts["<id>"].(string)

	result := db.Where("id = ?", name).Delete(&message.Message{})

	if result.Error != nil {
		fmt.Fprintf(os.Stderr, "cannot remove message %s: %v\n", name, result.Error)
		os.Exit(1)
	}
}
