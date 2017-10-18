package helper

import "crypto/rand"

// Secret is the JWT signing key
var Secret = generateRandomBytes(30)

// generateRandomBytes generated a []byte containing n secure random numbers
func generateRandomBytes(n int) []byte {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		panic(err)
	}
	return b
}
