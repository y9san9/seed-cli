package session

import (
	"encoding/base64"
	"fmt"
	"github.com/atotto/clipboard"
	"os"
	"path"
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

	fileNameEncrypted, err := encrypt(key, []byte(info.Name()))
	fileName := base64.URLEncoding.EncodeToString(fileNameEncrypted)
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

	// File in the same directory, but mangled
	absolutePath := path.Join(
		path.Dir(text),
		fileName+".seed",
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
