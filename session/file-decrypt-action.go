package session

import (
	"encoding/base64"
	"io"
	"bytes"
	"bufio"
	"fmt"
	"github.com/atotto/clipboard"
	"os"
	"path"
	"strings"
)

func fileDecryptAction(
	name string,
	key []byte,
	text string,
) (bool, error) {
	if path.Ext(text) != ".seed" {
		return false, nil
	}
	info, err := os.Stat(text)
	if err != nil {
		return false, nil
	}
	if info.IsDir() {
		return false, nil
	}

	fmt.Println("Decrypting file...")

	file, err := os.Open(text)
	if err != nil {
		return false, err
	}
	defer file.Close()
	reader := bufio.NewReader(file)

	fileName, err := decryptFileName(reader, key)
	if err != nil {
		return false, err
	}
	absolutePath := path.Join(path.Dir(text), fileName)

	decrypted, err := decryptRest(reader, key)
	if err != nil {
		return false, err
	}

	err = os.WriteFile(absolutePath, decrypted, 0600)
	if err != nil {
		return false, err
	}

	err = clipboard.WriteAll(absolutePath)
	if err == nil {
		fmt.Println(name+":", absolutePath, "(copied)")
	} else {
		fmt.Println(name+":", absolutePath)
	}

	return true, nil
}

func decryptFileName(
	reader *bufio.Reader,
	key []byte,
) (string, error) {
	encrypted, err := reader.ReadBytes('\n')
	if err != nil {
		return "", err
	}
	encrypted = bytes.TrimRight(encrypted, "\n")
	encrypted, err = base64.URLEncoding.DecodeString(string(encrypted))

	decrypted, err := decrypt(key, encrypted)
	if err != nil {
		fmt.Printf("%w\n", err)
		return "", err
	}

	fileNameString := string(decrypted)
	originalExtension := path.Ext(fileNameString)
	if len(originalExtension) == 0 {
		originalExtension = ".decrypted"
	}
	originalName := strings.TrimSuffix(fileNameString, originalExtension)
	newName := originalName + ".seed" + originalExtension

	return newName, nil
}

func decryptRest(
	reader *bufio.Reader,
	key []byte,
) ([]byte, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	return decrypt(key, data)
}
