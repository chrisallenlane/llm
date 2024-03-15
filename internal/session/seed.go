package session

import (
	"fmt"
	"runtime"

	"gorm.io/gorm"
)

// Seed seeds the database with session data
func Seed(db *gorm.DB) error {

	var assistant = `
You are a helpful assistant.
The user will ask you various questions. Please answer them to the best of your abilities.
Where possible, include links and citations in an unordered list at the end of your response.
Prefer official, high-quality sources where available.`

	var codegen = `
You are a code-completion tool used by software engineers.
The user will ask questions regarding how to accomplish various tasks in code.
Respond ONLY with code snippets that answer the user's question.
Your responses will be copied and pasted directly into the user's editor.
Wrap responses in markdown code fences.`

	var shell = fmt.Sprintf(`
You are a system administrator.
The user will ask questions regarding how to accomplish a task on their %s system.
Return the shell commands necessary to accomplish the task.
Wrap shell commands in markdown code fences.`, runtime.GOOS)

	// define the seeds
	seeds := []Session{
		Session{
			Name: "assistant",
			Hint: assistant,
		},
		Session{
			Name: "codegen",
			Hint: codegen,
		},
		Session{
			Name: "shell",
			Hint: shell,
		},
	}

	// ensure that each seed has been created
	for _, seed := range seeds {
		// check if a session already exists by this name
		result := db.
			Where("name = ?", seed.Name).
			Limit(1).
			Find(&Session{})
		if result.Error != nil {
			return fmt.Errorf(
				"failed to query seed: %s, %v",
				seed.Name,
				result.Error,
			)
		}

		// if the session already exists, do nothing
		if result.RowsAffected == 1 {
			continue
		}

		// otherwise, create it
		result = db.Create(&seed)
		if result.Error != nil {
			return fmt.Errorf(
				"failed to insert seed: %s, %v",
				seed.Name,
				result.Error,
			)
		}
	}

	return nil
}
