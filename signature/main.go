package signature

import (
	"crypto/sha256"
	"encoding/base32"
	"strings"
)

func DisplayString(key []byte) string {
	signature := sha256.Sum256(key)
	signatureBase32 := base32.StdEncoding.EncodeToString(signature[:])
	signatureBase32 = strings.TrimRight(signatureBase32, "=")

	var signatureString string
	first := true
	for start := 0; start < len(signatureBase32); start += 8 {
		middle := start + 4
		end := start + 8
		if end > len(signatureBase32) {
			end = len(signatureBase32)
		}
		if !first {
			signatureString += " "
		} else {
			first = false
		}
		if end-start == 8 {
			signatureString += signatureBase32[start:middle]
			signatureString += "-"
			signatureString += signatureBase32[middle:end]
		} else {
			signatureString += signatureBase32[start:end]
		}
	}

	return signatureString
}
