package services

import (
	"fmt"
	"math/rand"
	"time"
)

func GenerateUniqueCustomerID() string {
	// Implement your logic to generate a unique customer ID (e.g., UUID, timestamp, etc.)
	// For example, you can use a combination of timestamp and random characters
	return fmt.Sprintf("%d%s", time.Now().UnixNano(), GetRandomString(4))
}

// Custom function to generate random characters (for demonstration purposes)
func GetRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}


