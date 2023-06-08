package utils

import (
	"math"
	"math/rand"
	"time"
)

func GenerateReceiveCode(totalEntries int) (string, error) {
	codeLength := int(math.Log10(float64(totalEntries))) + 2
	if codeLength < 6 {
		codeLength = 6
	}
	rand.New(rand.NewSource(time.Now().UnixNano()))
	code := make([]byte, codeLength)
	for i := 0; i < codeLength; i++ {
		code[i] = byte(rand.Intn(10) + '0')
	}
	return string(code), nil
}
