package plugin

import (
	"crypto/rand"
	"encoding/hex"
	"io"
)

func randStr(n int) string {
	b := make([]byte, n/2)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}
