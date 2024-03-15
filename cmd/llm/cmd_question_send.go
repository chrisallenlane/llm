package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/briandowns/spinner"
	"github.com/chrisallenlane/llm/internal/message"
	"github.com/chrisallenlane/llm/internal/session"
	"github.com/sashabaranov/go-openai"
	"gorm.io/gorm"
)

// cmdQuestionSend sends the question stack to GPT
func cmdQuestionSend(_ map[string]interface{}, sess session.Session, db *gorm.DB) {
	// initialize the role hint
	messageSystem := message.Message{
		Content:   sess.Hint,
		Date:      time.Now(),
		Name:      "", // NB: this is unused
		Role:      openai.ChatMessageRoleSystem,
		SessionID: sess.ID,
	}

	// pack the messages for delivery to the API
	var messages []openai.ChatCompletionMessage
	for _, m := range sess.Messages {
		messages = append(messages, m.Pack())
	}

	// Arrange the messages in the following order before sending them to the API:
	//
	//   1. The first `n` messages in the history
	//   2. The system (role) hint
	//   3. The user's message
	//
	// We do this because the API reads each message from oldest to newest for
	// establishing context when preparing a response.
	messages = append(messages, messageSystem.Pack())

	// start a spinner while we wait
	s := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
	s.Suffix = " Querying GPT"
	s.Start()

	// initialize a new completion request
	request := openai.ChatCompletionRequest{
		// TODO: use the specified model
		Model:    openai.GPT3Dot5Turbo,
		Messages: messages,
	}

	// call to the API
	client := openai.NewClient(os.Getenv("OPENAI_KEY"))
	completionResponse, err := client.CreateChatCompletion(context.Background(), request)

	// stop the spinner
	s.Stop()

	// handle API errors
	if err != nil {
		fmt.Fprintf(os.Stderr, "completion API errored: %v\n", err)
		os.Exit(1)
	}

	// unpack the response content
	response := completionResponse.Choices[0].Message
	msg := message.Unpack(response, sess.ID)

	// save the response unless it was a one-off question
	if err := db.Create(&msg).Error; err != nil {
		fmt.Fprintf(os.Stderr, "failed to save response message: %v\n", err)
		os.Exit(1)
	}

	// display the message response
	msg.Display()
}
