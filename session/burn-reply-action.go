package session

import (
	"crypto/ecdh"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/atotto/clipboard"
	"github.com/fxamacker/cbor/v2"
)

// todo: check if received burn-reply without burn-request
func burnReplyAction(
	input *inputState,
	decrypted []byte,
) (bool, error) {
	curve := ecdh.X25519()
	key := input.key

	var wrongConfirmation burnConfirmation
	err := cbor.Unmarshal(decrypted, &wrongConfirmation)
	if err == nil && wrongConfirmation.Type == "burn_confirmation" {
		fmt.Println("Received reply to self-burning session while not in session.")
		fmt.Println("Use :burn and try again.")
		return true, nil
	}

	var request burnRequest
	err = cbor.Unmarshal(decrypted, &request)
	if err != nil || request.Type != "burn_request" {
		return false, nil
	}

	requestPublicKey, err := curve.NewPublicKey(request.PublicKey)
	if err != nil {
		return false, err
	}
	privateKey, err := curve.GenerateKey(rand.Reader)
	if err != nil {
		return false, err
	}
	burnKey, err := privateKey.ECDH(requestPublicKey)
	if err != nil {
		return false, err
	}
	confirmationCbor, err := cbor.Marshal(burnConfirmation{
		Type:      "burn_confirmation",
		PublicKey: privateKey.PublicKey().Bytes(),
	})
	if err != nil {
		return false, err
	}
	confirmationEncrypted, err := encrypt(key, confirmationCbor)
	if err != nil {
		return false, err
	}
	confirmationString := base64.URLEncoding.EncodeToString(confirmationEncrypted)

	fmt.Println("SELF-BURNING SESSION")
	fmt.Println("You received a self-burning session request.")
	fmt.Println("All keys are stored IN MEMORY, so after restart there is NO WAY TO DECRYPT message contents again.")
	fmt.Println("Used to ensure safety of both parties.")

	fmt.Printf("Step 1. Burn confirmation: %s", confirmationString)

	err = clipboard.WriteAll(confirmationString)
	if err == nil {
		fmt.Println(" (copied)")
	} else {
		fmt.Println()
	}

	fmt.Printf("Step 2. Ask %s to decrypt burn confirmation\n", input.peer)
	fmt.Println("Self-burning session started!")
	input.burnKey = burnKey
	return true, nil
}
