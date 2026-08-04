package storage

import (
    "os"
    "path/filepath"
    "encoding/base64"
)

func ensureSeedDirectory() string {
    home, err := os.UserHomeDir()
    if err != nil {
        panic(err)
    }
    directory := filepath.Join(home, ".seed")
    error := os.MkdirAll(directory, 0700)
    if error != nil {
        panic(error)
    }
    return directory
}

func ensurePeersDirectory() string {
    directory := filepath.Join(
        ensureSeedDirectory(),
        "peers",
    )
    error := os.MkdirAll(directory, 0700)
    if error != nil {
        panic(error)
    }
    return directory
}

func peerPath(name string) string {
    return filepath.Join(ensurePeersDirectory(), name+".shared")
}

func HasPeer(name string) bool {
    _, err := os.Stat(peerPath(name))
    return err == nil
}

func SavePeer(
    name string,
    shared []byte,
) {
    path := peerPath(name)
    string := base64.RawStdEncoding.EncodeToString(shared)
    err := os.WriteFile(path, []byte(string), 0600)
    if err != nil {
        panic(err)
    }
}

func LoadPeer(name string) ([]byte, bool) {
    path := peerPath(name)
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, false
    }
    result, err := base64.RawStdEncoding.DecodeString(string(data))
    if err != nil {
        panic(err)
    }
    return result, true
}
