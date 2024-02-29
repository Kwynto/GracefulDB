package core

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"os"
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
func CheckFolder(patch, name string) bool {
	// This function is complete
	fullPath := fmt.Sprintf("%s%s", patch, name)
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

func Encode64(inStr string) string {
	return base64.StdEncoding.EncodeToString([]byte(inStr))
}

func Decode64(inB64 string) string {
	decodeData, err := base64.StdEncoding.DecodeString(inB64)
	if err != nil {
		return ""
	}
	return string(decodeData)
}
