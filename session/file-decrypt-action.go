package session

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/atotto/clipboard"
	"github.com/fxamacker/cbor/v2"
	"io"
	"os"
	"path"
	"strings"
)

func fileDecryptAction(
	input *inputState,
) (bool, error) {
	text := input.text
	key := input.key
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
	metadataBytes, err := readMetadataBytes(reader)
	if err != nil {
		return false, err
	}

	// set metadataBytes and reader
	reader, metadataBytes, burn, err := readBurnFile(reader, metadataBytes, input.burnKey)
	if err != nil {
		return false, err
	}
	if burn && input.burnKey == nil {
		fmt.Println("File was sent during self-burning session. Keys were destroyed :D")
		return true, nil
	}
	if !burn && input.burnKey != nil {
		fmt.Printf("%s left self-burning session, leaving as well...\n", input.peer)
		input.burnKey = nil
	}

	var metadata fileMetadata
	err = cbor.Unmarshal(metadataBytes, &metadata)
	if err != nil {
		return false, err
	}
	filename := transformFilename(metadata.Name)
	absolutePath := path.Join(path.Dir(text), filename)

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
		fmt.Println(input.peer+":", absolutePath, "(copied)")
	} else {
		fmt.Println(input.peer+":", absolutePath)
	}

	return true, nil
}

func readBurnFile(
	reader *bufio.Reader,
	metadataBytes []byte,
	burnKey []byte,
) (*bufio.Reader, []byte, bool, error) {
	var burn burnFileMetadata
	err := cbor.Unmarshal(metadataBytes, &burn)
	if err != nil || burn.Type != "burn" {
		return reader, metadataBytes, false, nil
	}
	if burnKey == nil {
		return reader, metadataBytes, true, nil
	}

	content, err := io.ReadAll(reader)
	if err != nil {
		return nil, nil, false, err
	}

	decrypted, err := decrypt(burnKey, content)
	if err != nil {
		return nil, nil, false, err
	}

	resultReader := bufio.NewReader(bytes.NewReader(decrypted))
	resultMetadataBytes, err := readMetadataBytes(resultReader)
	if err != nil {
		return nil, nil, false, err
	}

	return resultReader, resultMetadataBytes, true, nil
}

func readMetadataBytes(reader *bufio.Reader) (metadataBytes []byte, err error) {
	var length int32
	err = binary.Read(reader, binary.LittleEndian, &length)
	if err != nil {
		return
	}
	metadataBytes = make([]byte, length)
	read, err := reader.Read(metadataBytes)
	if err != nil || int32(read) != length {
		return
	}
	return
}

func transformFilename(name string) string {
	originalExtension := path.Ext(name)
	if len(originalExtension) == 0 {
		originalExtension = ".decrypted"
	}
	originalName := strings.TrimSuffix(name, originalExtension)
	return originalName + ".seed" + originalExtension
}
