package session

import (
	"encoding/base64"
	"fmt"
	"github.com/fxamacker/cbor/v2"
)

func decryptAction(
	input *inputState,
) (bool, error) {
	text := input.text

	fromBase64, err := base64.URLEncoding.DecodeString(text)
	if err != nil {
		return false, nil
	}
	decrypted, err := decrypt(input.key, fromBase64)
	if err != nil {
		return false, nil
	}

	ok, err := burnReplyAction(input, decrypted)
	if err != nil {
		fmt.Printf("Received burn request, but couldn't answer it :( %v\n", err)
		return true, nil
	}
	if ok {
		return true, nil
	}

	var message message
	err = cbor.Unmarshal([]byte(decrypted), &message)
	if err != nil {
		return false, err
	}

	if message.Burn && input.burnKey == nil {
		fmt.Println("Message was sent during self-burning session. Keys are destroyed :D")
		return true, nil
	}
	if !message.Burn && input.burnKey != nil {
		fmt.Printf("%s left self-burning session, leaving as well...\n", input.peer)
		input.burnKey = nil
	}
	var payload []byte
	if input.burnKey != nil {
		payload, err = decrypt(input.burnKey, message.Payload)
		if err != nil {
			return false, err
		}
	} else {
		payload = message.Payload
	}

	fmt.Print(input.peer)
	if input.burnKey != nil {
		fmt.Print(" (BURN)")
	}
	fmt.Println(":", string(payload))
	return true, nil
}
