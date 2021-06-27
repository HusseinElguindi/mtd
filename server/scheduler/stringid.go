package scheduler

import (
	"crypto/rand"
	"encoding/hex"
	"strconv"
)

const idLen = 12

func truncStr(s string, trimTo int) string {
	if len(s) < trimTo {
		trimTo = len(s)
	}
	return s[:trimTo]
}

// https://github.com/moby/moby/blob/20.10/pkg/stringid/stringid.go
func GenerateID() string {
	b := make([]byte, 32)
	for {
		if _, err := rand.Read(b); err != nil {
			panic(err) // This shouldn't happen
		}
		id := hex.EncodeToString(b)
		// if we try to parse the truncated for as an int and we don't have
		// an error then the value is all numeric and causes issues when
		// used as a hostname. ref #3869
		if _, err := strconv.ParseInt(truncStr(id, idLen), 10, 64); err == nil {
			continue
		}
		return id
	}
}
