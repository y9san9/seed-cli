package session

import (
	"encoding/hex"
	"encoding/base64"
	"fmt"
	"github.com/atotto/clipboard"
	"os"
	"path"
	"crypto/sha256"
	"bytes"
)

func fileEncryptAction(
	key []byte,
	text string,
) (bool, error) {
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
	encryptedFileName, err := encryptFileName(key, info.Name())
	if err != nil {
		return false, err
	}

	data, err := os.ReadFile(text)
	if err != nil {
		return false, err
	}

	encrypted, err := encrypt(key, data)
	if err != nil {
		return false, err
	}

	var buffer bytes.Buffer
	buffer.WriteString(encryptedFileName+"\n")
	buffer.Write(encrypted)

	absolutePath := path.Join(
		path.Dir(text),
		hashedFileName+".seed",
	)

	if err := os.WriteFile(absolutePath, buffer.Bytes(), 0600); err != nil {
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

func encryptFileName(key []byte, name string) (string, error) {
	fileNameEncrypted, err := encrypt(key, []byte(name))
	if err != nil {
		return "", err
	}
	fileName := base64.URLEncoding.EncodeToString(fileNameEncrypted)
	return fileName, nil
}
