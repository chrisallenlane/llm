package input

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
)

// Editor reads input from the user's $EDITOR
func Editor(text string, _ map[string]interface{}) (string, error) {

	// create a temporary file
	tempFile, err := ioutil.TempFile("", "tempfile_*.txt")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	// write the default text to the temp file
	_, err = tempFile.WriteString(text)
	if err != nil {
		return "", fmt.Errorf("failed to write default text to file: %v", err)
	}

	// Get the user's preferred editor from the environment or use a default
	// TODO: support VISUAL
	editor := os.Getenv("EDITOR")
	if editor == "" {
		// TODO: make this configurable
		editor = "sensible-editor"
	}

	// Run the editor with the temporary file as an argument
	cmd := exec.Command(editor, tempFile.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		return "", fmt.Errorf("failed to run editor: %v", err)
	}

	// Read the contents of the temporary file
	contents, err := ioutil.ReadFile(tempFile.Name())
	if err != nil {
		return "", fmt.Errorf("failed to read temp file: %v", err)
	}

	return string(contents), nil
}
