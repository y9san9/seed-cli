package session

import (
	"bufio"
	"fmt"
	"github.com/atotto/clipboard"
	"os"
	"seed/storage"
	"strings"
)

type message struct {
	Version int    `json:"version"`
	Text    string `json:"text"`
}

func Run(name string) {
	scanner := bufio.NewScanner(os.Stdin)
	key, success := storage.LoadPeer(name)
	if !success {
		fmt.Printf("There is no session with name '%s'.\n", name)
		fmt.Println("To check available peers, visit ~/.seed/peers")
		return
	}

	fmt.Printf("Session with %s has started.\n", name)
	fmt.Println()
	fmt.Println("TIPS:")
	fmt.Println("• Messages will be decoded or encoded automatically.")
	fmt.Println("• Use '\\' for newlines.")
	fmt.Println("• Enter file path to encrypt/decrypt whole file.")
	fmt.Println("• Text is copied to clipboard when encoded.")
	fmt.Println("• Leave message empty to paste from clipboard.")
	fmt.Println("• Use :e to edit prompt with a proper editor.")
	fmt.Println()
	fmt.Println("Paste text in the text field below:")

	for {
		text := scanText(scanner)

		newText, ok, err := editorAction(text)
		if err != nil {
			fmt.Println("Couldn't use editor :(")
			continue
		}
		if ok {
			text = newText
		}

		if len(text) == 0 {
			continue
		}

		ok, err = fileEncryptAction(key, text)
		if err != nil {
			fmt.Println("Message is poorely formatted.")
			continue
		}
		if ok {
			continue
		}
		ok, err = fileDecryptAction(name, key, text)
		if err != nil {
			fmt.Println("Message is poorely formatted.")
			continue
		}
		if ok {
			continue
		}
		ok, err = decryptAction(name, key, text)
		if err != nil {
			fmt.Println("Message is poorely formatted.")
			continue
		}
		if ok {
			continue
		}
		err = encryptAction(key, text)
		if err != nil {
			fmt.Printf("Can't encrypt this message :(\n", err)
			continue
		}
	}

	fmt.Println("Session finished, bye :D")
}

func scanText(scanner *bufio.Scanner) string {
	var text string
	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}
		text += strings.TrimSpace(scanner.Text())
		if !strings.HasSuffix(text, "\\") {
			break
		}
		text = strings.TrimSuffix(text, "\\")
		text += "\n"
	}
	text = strings.TrimSpace(text)
	if len(text) == 0 {
		var err error
		text, err = clipboard.ReadAll()
		if err != nil {
			fmt.Println("Cannot read clipboard!")
		} else if len(text) > 0 {
			fmt.Printf("Clipboard read (%d bytes)\n", len([]byte(text)))
		}
	}
	return strings.TrimSpace(text)
}
