package session

import (
	"encoding/base64"
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
	fileName := strings.TrimSuffix(info.Name(), ".seed")

	fileNameEncrypted, err := base64.URLEncoding.DecodeString(fileName)
	if err != nil {
		return false, err
	}

	fileNameDecrypted, err := decrypt(key, fileNameEncrypted)
	if err != nil {
		return false, err
	}
	fileNameString := string(fileNameDecrypted)

	originalExtension := path.Ext(fileNameString)
	if len(originalExtension) == 0 {
		originalExtension = ".decrypted"
	}
	originalName := strings.TrimSuffix(fileNameString, originalExtension)
	newName := originalName + ".seed" + originalExtension
	absolutePath := path.Join(
		path.Dir(text),
		newName,
	)

	data, err := os.ReadFile(text)
	if err != nil {
		return false, err
	}
	decrypted, err := decrypt(key, data)
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
