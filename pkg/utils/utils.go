package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/sefazi/machine-learning-disease-diagnosis-apps-backend/internal/config"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes)
}

func Uint64ToBytes(n uint64) []byte {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, n)
	return buf
}

func BytesToUint64(b []byte) uint64 {
	return binary.BigEndian.Uint64(b)
}

func EncryptUint64(n uint64) string {
	cfg := config.GetConfig()
	plaintext := Uint64ToBytes(n)

	block, _ := aes.NewCipher(cfg.Config.EncryptionKey)

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return ""
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	// Base64 URL-safe encoding (no slashes, plus, or padding)
	return base64.RawURLEncoding.EncodeToString(ciphertext)
}

func DecryptToUint64(encoded string) uint64 {
	cfg := config.GetConfig()

	ciphertext, err := base64.RawURLEncoding.DecodeString(encoded)
	if err != nil || len(ciphertext) < aes.BlockSize {
		return 0
	}

	block, _ := aes.NewCipher(cfg.Config.EncryptionKey)
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	return BytesToUint64(ciphertext)
}

func GenerateRegisterNumber(name string) string {
	parts := strings.Fields(name)
	initials := ""

	for _, word := range parts {
		initials += strings.ToUpper(string(word[0]))
	}

	timestamp := time.Now().Format("20060102-150405")

	if initials == "" {
		return timestamp
	}

	return fmt.Sprintf("%s-%s", initials, timestamp)
}
