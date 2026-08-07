package session

import (
	"encoding/base64"
	"fmt"
	"github.com/atotto/clipboard"
	"github.com/fxamacker/cbor/v2"
)

func encryptAction(
	input *inputState,
) error {
	var payload []byte
	var err error

	if input.burnKey != nil {
		payload, err = encrypt(input.burnKey, []byte(input.text))
		if err != nil {
			return err
		}
	} else {
		payload = []byte(input.text)
	}

	message := message{
		Version: 0,
		Payload: payload,
		Burn:    input.burnKey != nil,
	}

	cbor, err := cbor.Marshal(message)
	if err != nil {
		return err
	}

	encrypted, err := encrypt(input.key, cbor)
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
