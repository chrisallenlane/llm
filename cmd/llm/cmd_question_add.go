package main

import (
	"fmt"
	"os"
	"time"

	"github.com/chrisallenlane/llm/internal/input"
	"github.com/chrisallenlane/llm/internal/message"
	"github.com/chrisallenlane/llm/internal/session"
	"github.com/sashabaranov/go-openai"
	"gorm.io/gorm"
)

// cmdQuestionAdd stages a question to ask
func cmdQuestionAdd(opts map[string]interface{}, sess session.Session, db *gorm.DB) {
	// read the user's question
	question, err := input.Read("<msg>", opts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read input: %v\n", err)
		os.Exit(1)
	}

	// abort if the question is empty
	if question == "" {
		fmt.Fprintf(os.Stderr, "Aborted.\n")
		os.Exit(1)
	}

	// initialize the user's message
	messageUser := message.Message{
		Content:   question,
		Date:      time.Now(),
		Name:      "", // NB: this is unused
		Role:      openai.ChatMessageRoleUser,
		SessionID: sess.ID,
	}

	// save the user's question to the database
	if err := db.Create(&messageUser).Error; err != nil {
		fmt.Fprintf(os.Stderr, "failed to save user message: %v\n", err)
		os.Exit(1)
	}
}
