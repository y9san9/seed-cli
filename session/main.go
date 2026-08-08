package session

import (
	"bufio"
	"fmt"
	"os"
	"seed/signature"
	"seed/storage"
)

func Run(args []string) {
	if len(args) < 3 {
		fmt.Println("Usage: seed session <peer>")
		fmt.Println("To check available peers, visit ~/.seed/peers")
		return
	}
	var input inputState
	input.peer = args[2]

	scanner := bufio.NewScanner(os.Stdin)
	key, success := storage.LoadPeer(input.peer)
	if !success {
		fmt.Printf("There is no peer with name '%s'.\n", input.peer)
		fmt.Println("To check available peers, visit ~/.seed/peers")
		return
	}
	input.key = key

	fmt.Printf("Session with %s has started.\n", input.peer)
	fmt.Println()
	fmt.Println("PUBLIC SIGNATURE:", signature.DisplayString(input.key), "(SHA-256)")
	fmt.Println()
	fmt.Println("If both parties see the same text here, it is GUARANTEED that no one can read messages other than you.")
	fmt.Println("If keys were MODIFIED by someone in the middle, signature will differ.")
	fmt.Println("Suggested method of verification is a DIFFERENT untrusted channel, ideally over a VOICE CALL.")
	fmt.Println()
	fmt.Println("TIPS:")
	fmt.Println("• Messages will be decoded or encoded automatically.")
	fmt.Println("• Use '\\' for newlines.")
	fmt.Println("• Enter file path to encrypt/decrypt whole file.")
	fmt.Println("• Text is copied to clipboard when encoded.")
	fmt.Println("• Leave message empty to paste from clipboard.")
	fmt.Println("• Use :e to edit prompt with a proper editor.")
	fmt.Println("• Use :burn to create a self-burning session.")
	fmt.Println()
	fmt.Println("Paste text in the text field below:")

	for {
		input.text = scanText(input.burnKey != nil, scanner)

		ok, err := burnRequestAction(&input)
		if err != nil {
			fmt.Printf("Couldn't start burn session. %+v\n", err)
			continue
		}
		if ok {
			continue
		}

		err = editorAction(&input)
		if err != nil {
			fmt.Println("Couldn't use editor :(", err)
			continue
		}

		if len(input.text) == 0 {
			continue
		}

		ok, err = fileEncryptAction(&input)
		if err != nil {
			fmt.Printf("Message is poorely formatted. %#v\n", err)
			continue
		}
		if ok {
			continue
		}
		ok, err = fileDecryptAction(&input)
		if err != nil {
			fmt.Printf("Message is poorely formatted. %+v\n", err)
			continue
		}
		if ok {
			continue
		}
		ok, err = decryptAction(&input)
		if err != nil {
			fmt.Println("Message is poorely formatted.", err)
			continue
		}
		if ok {
			continue
		}
		err = encryptAction(&input)
		if err != nil {
			fmt.Printf("Can't encrypt this message :( %+w\n", err)
			continue
		}
	}

	fmt.Println("Session finished, bye :D")
}
