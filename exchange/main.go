package exchange

import (
	"crypto/ecdh"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"seed/storage"
	"github.com/atotto/clipboard"
)

const keySize = 32
const incorrectKeyMessage = "Paste the correct key, try again: "

func Run() {
	curve := ecdh.X25519()

	fmt.Println("Let's start an exchanging mechanism...")
	fmt.Print("Peer name: ")

	var name string
	for {
		fmt.Scanln(&name)

		if storage.HasPeer(name) {
			fmt.Printf("Peer with name %s already exists\n", name)
			continue
		}

		break
	}

	private, err := curve.GenerateKey(rand.Reader)
	if err != nil {
		panic(err)
	}

	public := base64.URLEncoding.EncodeToString(private.PublicKey().Bytes())

	fmt.Printf("=========== Initializing %s's Peer ===========\n", name)

	fmt.Printf("Step 1. Send this over untrusted channel: %s", public)

	err = clipboard.WriteAll(public)
	if err == nil {
		fmt.Println(" (copied)")
	} else {
		fmt.Println()
	}

	fmt.Printf("Step 2. Paste %s's response: ", name)

	var sharedResult []byte
	for {
		var peer string
		fmt.Scanln(&peer)

		peerBytes, err := base64.URLEncoding.DecodeString(peer)
		if err != nil {
			fmt.Print(incorrectKeyMessage)
			continue
		}

		peerPublic, err := curve.NewPublicKey(peerBytes)
		if err != nil {
			fmt.Print(incorrectKeyMessage)
			continue
		}

		shared, err := private.ECDH(peerPublic)
		if err != nil {
			fmt.Print(incorrectKeyMessage)
			continue
		}

		sharedResult = shared
		break
	}

	storage.SavePeer(name, sharedResult)

	fmt.Println("=========== Shared key generated ===========")
	fmt.Printf("Now you can encode/decode messages using 'seed session %s' command\n", name)
}
