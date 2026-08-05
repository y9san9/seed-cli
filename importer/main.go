package importer

import (
	"crypto/ecdh"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
    "seed/storage"
	"github.com/atotto/clipboard"
)

type export struct {
    PublicKey string `json:"public_key"`
    Name      string `json:"name"`
    Payload   string `json:"payload"`
}

const incorrectKeyMessage = "Enter a valid secured key, try again: "

func Run() {
	curve := ecdh.X25519()

	privateKey, err := curve.GenerateKey(rand.Reader)
	if err != nil {
		panic(err)
	}
    publicKeyString := base64.URLEncoding.EncodeToString(
        privateKey.PublicKey().Bytes(),
    )

    fmt.Println("Ready to receive peer securely.")
    fmt.Println("The actual shared key WILL NEVER be exposed.")
    fmt.Println("You can share single-time key in untrusted environments.")
    fmt.Println("Run 'seed export' on the other machine to get further instructions.")
    fmt.Print("Step 1. Use single-time key: ", publicKeyString)

	err = clipboard.WriteAll(publicKeyString)
	if err == nil {
		fmt.Println(" (copied)")
	} else {
		fmt.Println()
	}

    fmt.Print("Step 2. Secured key: ")

	for {
		var exportString string
		fmt.Scanln(&exportString)

        exportBytes, err := base64.URLEncoding.DecodeString(exportString)
        if err != nil {
			fmt.Print(incorrectKeyMessage)
			continue
        }

        var export export
        err = json.Unmarshal(exportBytes, &export)
        if err != nil {
			fmt.Print(incorrectKeyMessage)
			continue
        }

        fmt.Println(export)

		publicKeyBytes, err := base64.URLEncoding.DecodeString(export.PublicKey)
		if err != nil {
			fmt.Print(incorrectKeyMessage)
			continue
		}

		publicKey, err := curve.NewPublicKey(publicKeyBytes)
		if err != nil {
			fmt.Print(incorrectKeyMessage)
			continue
		}

        singleTimeSecret, err := privateKey.ECDH(publicKey)
		if err != nil {
			fmt.Print(incorrectKeyMessage)
			continue
		}

        nameBytes, err := base64.URLEncoding.DecodeString(export.Name)
		if err != nil {
			fmt.Print(incorrectKeyMessage)
			continue
		}

        nameBytes, err = decrypt(singleTimeSecret, nameBytes)
		if err != nil {
			fmt.Print(incorrectKeyMessage)
			continue
		}

        name := string(nameBytes)

        if storage.HasPeer(name) {
            fmt.Printf("Peer with name %s already exists, try again: ", name)
            break
        }

        staticSecret, err := base64.URLEncoding.DecodeString(export.Payload)
        if err != nil {
			fmt.Print(incorrectKeyMessage)
			continue
        }

        staticSecret, err = decrypt(singleTimeSecret, staticSecret)
        if err != nil {
			fmt.Print(incorrectKeyMessage)
			continue
        }

        storage.SavePeer(name, staticSecret)

        fmt.Printf("=========== Peer %s shared ===========\n", name)
        fmt.Printf("Now you can encode/decode messages using 'seed session %s' command\n", name)
		break
	}
}
