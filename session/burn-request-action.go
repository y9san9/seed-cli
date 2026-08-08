package session

import (
	"crypto/ecdh"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/atotto/clipboard"
	"github.com/fxamacker/cbor/v2"
)

func burnRequestAction(input *inputState) (bool, error) {
	curve := ecdh.X25519()
	if input.text != ":burn" {
		return false, nil
	}
	privateKey, err := curve.GenerateKey(rand.Reader)
	if err != nil {
		return false, err
	}

	burnRequestCbor, err := cbor.Marshal(burnRequest{
		Type:      "burn_request",
		PublicKey: privateKey.PublicKey().Bytes(),
	})
	if err != nil {
		return false, err
	}
	burnRequestEncrypted, err := encrypt(input.key, []byte(burnRequestCbor))
	if err != nil {
		return false, err
	}
	burnRequestString := base64.URLEncoding.EncodeToString(burnRequestEncrypted)

	fmt.Println("SELF-BURNING SESSION")
	fmt.Println("To create a self-burning session we first need to do some math magic.")
	fmt.Println("All keys for that session are stored IN MEMORY, so after restart there is NO WAY TO DECRYPT message contents again.")
	fmt.Println("Used to ensure safety of both parties.")

	fmt.Printf("Step 1. Burn request: %s", burnRequestString)
	err = clipboard.WriteAll(burnRequestString)
	if err == nil {
		fmt.Println(" (copied)")
	} else {
		fmt.Println()
	}

	fmt.Printf("Step 2. Ask %s to decrypt burn request while in session.\n", input.peer)

	fmt.Print("Step 3. Burn confirmation: ")

	var burnKey []byte
	for {
		var confirmationString string
		fmt.Scanln(&confirmationString)
		confirmationBytes, err := base64.URLEncoding.DecodeString(confirmationString)
		if err != nil {
			fmt.Print("Enter a valid confirmation, try again: ")
			continue
		}
		confirmationDecrypted, err := decrypt(input.key, confirmationBytes)
		if err != nil {
			fmt.Print("Enter a valid confirmation, try again: ")
			continue
		}
		var confirmation burnConfirmation
		err = cbor.Unmarshal(confirmationDecrypted, &confirmation)
		if err != nil {
			fmt.Print("Enter a valid confirmation, try again: ")
			continue
		}
		confirmationPublicKey, err := curve.NewPublicKey(confirmation.PublicKey)
		if err != nil {
			fmt.Print("Enter a valid confirmation, try again: ")
			continue
		}
		burnKey, err = privateKey.ECDH(confirmationPublicKey)
		if err != nil {
			fmt.Print("Enter a valid confirmation, try again: ")
			continue
		}
		break
	}

	fmt.Println("Self-burning session has started!")
	input.burn.key = burnKey
	return true, nil
}
