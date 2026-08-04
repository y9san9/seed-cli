package exchange

import (
    "fmt"
    "crypto/rand"
    "crypto/ecdh"
    "encoding/base64"
    "seed/storage"
)


const keySize = 32
const incorrectKeyMessage = "Paste the correct key, try again: "

func Run() {
    curve := ecdh.X25519()

    fmt.Println("Welcome to Seed Toolkit!")
    fmt.Println("Let's start an exchanging mechanism...")
    fmt.Print("Peer name: ")

    var name string
    fmt.Scanln(&name)

    private, err := curve.GenerateKey(rand.Reader)
    if err != nil {
        panic(err)
    }

    public := base64.RawStdEncoding.EncodeToString(private.PublicKey().Bytes())

    fmt.Printf("=========== Initializing %s's Peer ===========\n", name)
    fmt.Printf("Step 1. Send this over untrusted channel: %s\n", public)
    fmt.Printf("Step 2. Paste %s's response: ", name)

    var sharedResult []byte
    for {
        var peer string
        fmt.Scanln(&peer)

        peerBytes, err := base64.RawStdEncoding.DecodeString(peer)
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
