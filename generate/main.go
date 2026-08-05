package generate

import (
    "fmt"
    "crypto/rand"
    "crypto/ecdh"
    "seed/storage"
)

func Run() {
	curve := ecdh.X25519()

    fmt.Println("Local peer generation.")
    fmt.Println("Useful for latter distribution with 'seed export'.")
    fmt.Println("Main way to create circles where groups of people share the same key.")

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

    key, err := curve.GenerateKey(rand.Reader)
    if err != nil {
        panic(err)
    }
    storage.SavePeer(name, key.Bytes())

	fmt.Println("=========== Private key generated ===========")
}
