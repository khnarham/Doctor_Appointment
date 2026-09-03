package Infra
import (
	"crypto/rand"
	"encoding/hex"
)
type CryptoCodeGenerator struct{}

func NewCryptoCodeGenerator() *CryptoCodeGenerator {
	return &CryptoCodeGenerator{}
}

func (g *CryptoCodeGenerator) GenerateToken() string {
	bytes := make([]byte, 32) 
	if _, err := rand.Read(bytes); err != nil {
		return ""
	}
	return hex.EncodeToString(bytes)
}