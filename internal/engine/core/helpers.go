package core

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"
)

// Name generation
func GenerateName() string {
	// This function is complete
	b := make([]byte, 16)
	rand.Read(b)

	return fmt.Sprintf("%x", b)
}

// Revision generation
func GenerateRev() string {
	// This function is complete
	b := make([]byte, 4)
	rand.Read(b)

	return fmt.Sprintf("%x", b)
}

// Checking the folder name
func CheckFolder(path, name string) bool {
	// This function is complete
	// fullPath := fmt.Sprintf("%s%s", patch, name)
	fullPath := filepath.Join(path, name)
	dir, err := os.Stat(fullPath)
	if err != nil {
		return false
	}

	return dir.IsDir()
}

func Uint64ToBinary(i uint64) []byte {
	bs := make([]byte, 8)
	binary.BigEndian.PutUint64(bs, i)
	return bs
}

func BinaryToUint64(b []byte) uint64 {
	i := binary.BigEndian.Uint64(b)
	return i
}

func Encode64(input string) string { // input - ordinary string
	return base64.StdEncoding.EncodeToString([]byte(input))
}

func Decode64(input string) string { // input - base64 string
	decodeData, err := base64.StdEncoding.DecodeString(input)
	if err != nil {
		return ""
	}
	return string(decodeData)
}

func intPow(acc, base uint64, exponent uint8) uint64 {
	// This function is complete
	if exponent == 1 {
		return acc
	}
	return intPow(acc*base, base, exponent-1)
}

func Pow(base uint64, exponent uint8) uint64 {
	// This function is complete
	if exponent == 0 {
		return 1
	}
	if exponent == 1 {
		return base
	}
	return intPow(base*base, base, exponent-1)
}
