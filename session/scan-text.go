package session

import (
	"bufio"
	"fmt"
	"github.com/atotto/clipboard"
	"strings"
)

func scanText(
	burn bool,
	scanner *bufio.Scanner,
) string {
	var text string
	for {
		if burn {
			fmt.Print("(BURN) ")
		}
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
