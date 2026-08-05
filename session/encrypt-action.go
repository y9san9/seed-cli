package session

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/atotto/clipboard"
)

func encryptAction(
	key []byte,
	text string,
) error {
	message := message{
		Version: 0,
		Text:    text,
	}
	json, err := json.Marshal(message)
	if err != nil {
		return err
	}

	encrypted, err := encrypt(key, json)
	if err != nil {
		return err
	}

	base64 := base64.URLEncoding.EncodeToString(encrypted)

	err = clipboard.WriteAll(base64)
	if err == nil {
		fmt.Println("You:", base64, "(copied)")
	} else {
		fmt.Println("You:", base64)
	}
	return nil
}
