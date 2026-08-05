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

	data, err := os.ReadFile(text)
	if err != nil {
		return false, err
	}
	decrypted, err := decrypt(key, data)
	if err != nil {
		return false, err
	}
	reader := bufio.NewReader(bytes.NewReader(decrypted))

	fileName, err := readFileName(reader, key)
	if err != nil {
		return false, err
	}
	absolutePath := path.Join(path.Dir(text), fileName)

	content, err := io.ReadAll(reader)
	if err != nil {
		return false, err
	}

	err = os.WriteFile(absolutePath, content, 0600)
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

func readFileName(
	reader *bufio.Reader,
	key []byte,
) (string, error) {
	data, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	data = strings.TrimRight(data, "\n")
	decodedData, err := base64.URLEncoding.DecodeString(string(data))
	if err != nil {
		return "", err
	}
	string := string(decodedData)

	originalExtension := path.Ext(string)
	if len(originalExtension) == 0 {
		originalExtension = ".decrypted"
	}
	originalName := strings.TrimSuffix(string, originalExtension)
	newName := originalName + ".seed" + originalExtension

	return newName, nil
}
