package encryption_test

import (
	"encoding/hex"
	"math/rand"
	"techytechster/digitaldexterity/internal/encryption"
	"testing"
	"time"
)

func TestEncryption(t *testing.T) {
	t.Parallel()
	rand.Seed(time.Now().UnixNano())

	encrypted := encryption.Encrypt("AS3CUR3P455W0RD")

	if len(encrypted) == 0 {
		t.FailNow()
	}

	t.Log(hex.EncodeToString(encrypted))
}
