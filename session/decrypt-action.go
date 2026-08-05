package session

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
)

func decryptAction(
	name string,
	key []byte,
	text string,
) (bool, error) {
	fromBase64, err := base64.URLEncoding.DecodeString(text)
	if err != nil {
		return false, nil
	}
	decrypted, err := decrypt(key, fromBase64)
	if err != nil {
		return false, nil
	}
	var message message
	err = json.Unmarshal([]byte(decrypted), &message)
	if err != nil {
		return false, err
	}
	fmt.Printf("%s: %s\n", name, message.Text)
	return true, nil
}
