package session

import (
    "os"
    "os/exec"
    "fmt"
    "strings"
)

func editorAction(text string) (string, bool, error) {
    if text != ":e" {
        return "", false, nil
    }

    tempFile, err := os.CreateTemp("", "seed-edit-*.txt")
    if err != nil {
        return "", false, err
    }
    defer os.Remove(tempFile.Name())
    if err := tempFile.Close(); err != nil {
        return "", false, err
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
        return "", false, err
    }

    data, err := os.ReadFile(tempFile.Name())
    if err != nil {
        return "", false, err
    }

    string := strings.TrimSpace(string(data))

    if len(string) == 0 {
        fmt.Println("Editor returned no text.")
    } else {
        fmt.Println("Editor:", string)
    }
    return string, true, nil
}
