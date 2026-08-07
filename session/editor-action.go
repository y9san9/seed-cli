package session

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func editorAction(input *inputState) error {
	if input.text != ":e" {
		return nil
	}

	tempFile, err := os.CreateTemp("", "seed-edit-*.txt")
	if err != nil {
		return err
	}
	defer os.Remove(tempFile.Name())
	if err := tempFile.Close(); err != nil {
		return err
	}

	editor := os.Getenv("VISUAL")
	if editor == "" {
		editor = os.Getenv("EDITOR")
	}
	if editor == "" {
		editor = "vim"
	}

	cmd := exec.Command(editor, tempFile.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return err
	}

	data, err := os.ReadFile(tempFile.Name())
	if err != nil {
		return err
	}

	input.text = strings.TrimSpace(string(data))

	if len(input.text) == 0 {
		fmt.Println("Editor returned no text.")
	} else {
		fmt.Println("Editor:", input.text)
	}

	return nil
}
