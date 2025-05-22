package utils

import (
	"crypto/rand"
	"math/big"
)

// GenerateRandomOTP generates a 6-digit random OTP
func GenerateRandomOTP() string {
	const charset = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	const length = 6

	otp := make([]byte, length)
	for i := 0; i < length; i++ {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		otp[i] = charset[n.Int64()]
	}

	return string(otp)
}
