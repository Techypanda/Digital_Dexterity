package encryption

import (
	"crypto/sha256"
)

func Encrypt(toEncrypt string) []byte {
	h := sha256.New()
	h.Write([]byte(toEncrypt))
	sha := h.Sum(nil)

	return sha
}
