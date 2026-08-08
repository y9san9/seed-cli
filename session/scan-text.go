package session

import (
	"bufio"
	"context"
	"fmt"
	"github.com/atotto/clipboard"
	"strings"
)

func scanText(
	ctx context.Context,
	burn bool,
	scanner *bufio.Scanner,
) (string, error) {
	channel := make(chan bool, 1)
	var text string
	for {
		if burn {
			fmt.Print("(BURN) ")
		}
		fmt.Print("> ")
		go func() {
			channel <- scanner.Scan()
		}()
		var scan bool
		select {
		case <-ctx.Done():
			return "", ctx.Err()
		case scanValue := <-channel:
			scan = scanValue
		}
		if !scan {
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
	return strings.TrimSpace(text), nil
}
