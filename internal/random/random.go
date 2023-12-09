package random

import (
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func init() {
	rand.New(rand.NewSource(time.Now().UnixNano()))

}

func RandInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

func RandString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)] // --> to get a random character from the alphabet
		sb.WriteByte(c)
	}

	return sb.String()
}
