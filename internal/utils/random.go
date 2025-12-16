package utils

import (
	"math/rand"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// GenerateShortCode cria uma string aleatória de 'n' caracteres.
func GenerateShortCode(n int) string {
	b := make([]byte, n)
	for i := range b {
		// O rand.Intn global agora já é aleatório por padrão
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
