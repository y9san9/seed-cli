package export

import (
	"crypto/ecdh"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/atotto/clipboard"
	"seed/storage"
)

const incorrectKeyMessage = "Paste the correct key, try again: "

type export struct {
	PublicKey string `json:"public_key"`
	Name      string `json:"name"`
	Payload   string `json:"payload"`
}

func Run(args []string) {
	curve := ecdh.X25519()

	if len(args) < 3 {
		fmt.Println("Usage: seed export <peer>")
		fmt.Println("To check available peers, visit ~/.seed/peers")
		return
	}

	name := args[2]
	staticSecret, success := storage.LoadPeer(name)
	if !success {
		fmt.Printf("There is no peer with name '%s'.\n", name)
		fmt.Println("To check available peers, visit ~/.seed/peers")
		return
	}
	privateKey, err := curve.GenerateKey(rand.Reader)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Sharing key for peer %s.\n", name)
	fmt.Println("The actual shared key WILL NEVER be exposed.")
	fmt.Println("You can share single-time key in untrusted environments.")
	fmt.Println("In order to achieve that, we first need a single-time key from your target machine.")
	fmt.Println("ON TARGET MACHINE, type 'seed import' and paste it below.")
	fmt.Print("Step 1. Single-time key: ")

	var singleTimeSecret []byte
	for {
		var publicKeyString string
		fmt.Scanln(&publicKeyString)

		publicKeyBytes, err := base64.URLEncoding.DecodeString(publicKeyString)
		if err != nil {
			fmt.Print(incorrectKeyMessage)
			continue
		}

		publicKey, err := curve.NewPublicKey(publicKeyBytes)
		if err != nil {
			fmt.Print(incorrectKeyMessage)
			continue
		}

		singleTimeSecret, err = privateKey.ECDH(publicKey)
		if err != nil {
			fmt.Print(incorrectKeyMessage)
			continue
		}

		break
	}

	encryptedPayload, err := encrypt(singleTimeSecret, staticSecret)
	if err != nil {
		panic(err)
	}

	encryptedName, err := encrypt(singleTimeSecret, []byte(name))
	if err != nil {
		panic(err)
	}

	json, err := json.Marshal(export{
		PublicKey: base64.URLEncoding.EncodeToString(
			privateKey.PublicKey().Bytes(),
		),
		Name:    base64.URLEncoding.EncodeToString(encryptedName),
		Payload: base64.URLEncoding.EncodeToString(encryptedPayload),
	})
	if err != nil {
		panic(err)
	}

	stringPayload := base64.URLEncoding.EncodeToString([]byte(json))

	fmt.Print("Step 2. Import secured key: ", stringPayload)

	err = clipboard.WriteAll(stringPayload)
	if err == nil {
		fmt.Println(" (copied)")
	} else {
		fmt.Println()
	}

	fmt.Println("Next steps are completed on the other device.")
}
