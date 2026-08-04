package session

import (
    "fmt"
    "encoding/json"
    "crypto/aes"
    "crypto/cipher"
    "seed/storage"
    "crypto/rand"
    "encoding/base64"
    "io"
    "errors"
    "bufio"
    "os"
)

type message struct {
    Version int `json:"version"`
    Text string `json:"text"`
}

func Run(name string) {
    scanner := bufio.NewScanner(os.Stdin)
    fmt.Println("=========== Welcome to Seed Toolkit ===========")

    key, success := storage.LoadPeer(name)
    if !success {
        fmt.Printf("There is no session with name '%s'.\n", name)
        return
    }

    fmt.Printf("Session with %s has started.\n", name)
    fmt.Println()
    fmt.Println("TIPS:")
    fmt.Println("• Messages will be decoded or encoded automatically.")
    fmt.Println("• Use '\\' for newlines.")
    // fmt.Println("• Text is copied to clipboard when encoded.")
    // fmt.Println("• Leave message empty to paste from clipboard.")
    fmt.Println()
    fmt.Println("Paste text in the text field below:")

    for{
        fmt.Print("> ")
        if !scanner.Scan() {
            break
        }
        text := scanner.Text()

        decrypted, err := decrypt(key, text)
        if err == nil {
            var message message
            err := json.Unmarshal([]byte(decrypted), &message)
            if err != nil {
                fmt.Println("Message is poorely formatted.")
                continue
            }
            fmt.Printf("%s: %s\n", name, message.Text)
            continue
        }

        message := message {
            Version: 0,
            Text: text,
        }
        json, err := json.Marshal(message)
        if err != nil {
            fmt.Printf("Can't encode this message :(\n", err)
            continue
        }
        encrypted, err := encrypt(key, string(json))
        if err != nil {
            fmt.Printf("Can't encode this message :(\n", err)
            continue
        }
        fmt.Println("You:", encrypted)
    }

    fmt.Println("Session finished, bye :D")
}

func decrypt(
    key []byte,
    payload string,
) (string, error) {
    encrypted, err := base64.RawStdEncoding.DecodeString(payload)
    if err != nil {
        return "", err
    }
    block, err := aes.NewCipher(key)
    if err != nil {
        return "", err
    }
    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return "", err
    }
    if len(encrypted) <= gcm.NonceSize() {
        return "", errors.New("Improperly formatted message")
    }
    data, err := gcm.Open(nil, encrypted[:gcm.NonceSize()], encrypted[gcm.NonceSize():], nil)
    if err != nil {
        return "", err
    }
    return string(data), nil
}

func encrypt(
    key []byte,
    payload string,
) (string, error) {
    block, err := aes.NewCipher(key)
    if err != nil {
        return "", err
    }
    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return "", err
    }
    nonce := make([]byte, gcm.NonceSize())
    _, err = io.ReadFull(rand.Reader, nonce)
    if err != nil {
        return "", err
    }
    encrypted := gcm.Seal(nonce, nonce, []byte(payload), nil)
    base64 := base64.RawStdEncoding.EncodeToString(encrypted)
    return base64, nil
}
