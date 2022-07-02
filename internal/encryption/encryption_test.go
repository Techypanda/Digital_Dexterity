package encryption

import (
	"encoding/hex"
	"math/rand"
	"testing"
	"time"
)

func TestEncryption(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	encrypted := Encrypt("AS3CUR3P455W0RD")
	if len(encrypted) == 0 {
		t.FailNow()
	}
	t.Log(hex.EncodeToString(encrypted))
}
