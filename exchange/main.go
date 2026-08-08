package exchange

import (
	"crypto/ecdh"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/atotto/clipboard"
	"seed/signature"
	"seed/storage"
	"time"
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
			fmt.Printf("Peer with name %s already exists, try again: ", name)
			continue
		}

		break
	}

	private, err := curve.GenerateKey(rand.Reader)
	if err != nil {
		panic(err)
	}

	public := base64.URLEncoding.EncodeToString(private.PublicKey().Bytes())

	fmt.Println()
	fmt.Printf("=========== Initializing %s's Peer ===========\n", name)

	fmt.Printf("Step 1. Send this over untrusted channel: %s", public)

	err = clipboard.WriteAll(public)
	if err == nil {
		fmt.Println(" (copied)")
	} else {
		fmt.Println()
	}

	time.Sleep(5 * time.Second)

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

	fmt.Println()
	fmt.Println("=========== Shared key generated ===========")
	fmt.Println("PUBLIC SIGNATURE:", signature.DisplayString(sharedResult), "(SHA-256)")
	fmt.Println()
	time.Sleep(2 * time.Second)
	fmt.Println("If both parties see the same text here, it is GUARANTEED that no one can read messages other than you.")
	fmt.Println("If keys were MODIFIED by someone in the middle, signature will differ.")
	fmt.Println("Suggested method of verification is a DIFFERENT untrusted channel, ideally over a VOICE CALL.")
	fmt.Println()
	fmt.Printf("Now you can encode/decode messages using 'seed session %s' command\n", name)
}
