package session

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/atotto/clipboard"
	"github.com/fxamacker/cbor/v2"
	"os"
	"path"
)

func fileEncryptAction(input *inputState) (bool, error) {
	text := input.text
	key := input.key

	if path.Ext(text) == ".seed" {
		return false, nil
	}
	info, err := os.Stat(text)
	if err != nil {
		return false, nil
	}
	if info.IsDir() {
		return false, nil
	}

	fmt.Println("Encrypting file...")

	hashedFileName := hashFileName(info.Name())

	data, err := os.ReadFile(text)
	if err != nil {
		return false, err
	}

	metadata := fileMetadata{
		Type: "encrypted",
		Name: info.Name(),
	}
	metadataBytes, err := cbor.Marshal(metadata)
	if err != nil {
		return false, err
	}

	buffer := new(bytes.Buffer)
	length := int32(len(metadataBytes))
	if err := binary.Write(buffer, binary.LittleEndian, length); err != nil {
		return false, err
	}
	if _, err := buffer.Write(metadataBytes); err != nil {
		return false, err
	}
	if _, err := buffer.Write(data); err != nil {
		return false, err
	}

	if input.burnKey != nil {
		encrypted, err := encrypt(input.burnKey, buffer.Bytes())
		if err != nil {
			return false, err
		}
		metadata := burnFileMetadata{
			Type: "burn",
		}
		metadataBytes, err := cbor.Marshal(metadata)
		if err != nil {
			return false, err
		}
		buffer = new(bytes.Buffer)
		length := int32(len(metadataBytes))
		if err := binary.Write(buffer, binary.LittleEndian, length); err != nil {
			return false, err
		}
		if _, err := buffer.Write(metadataBytes); err != nil {
			return false, err
		}
		if _, err := buffer.Write(encrypted); err != nil {
			return false, err
		}
	}

	encrypted, err := encrypt(key, buffer.Bytes())
	if err != nil {
		return false, err
	}

	absolutePath := path.Join(
		path.Dir(text),
		hashedFileName+".seed",
	)

	if err := os.WriteFile(absolutePath, encrypted, 0600); err != nil {
		return false, err
	}

	err = clipboard.WriteAll(absolutePath)
	if err == nil {
		fmt.Println("You:", absolutePath, "(copied)")
	} else {
		fmt.Println("You:", absolutePath)
	}

	return true, nil
}

func hashFileName(name string) string {
	hash := sha256.Sum256([]byte(name))
	first8 := hash[:8]
	return hex.EncodeToString(first8)
}
