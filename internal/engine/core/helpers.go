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
func CheckFolder(sPath, sName string) bool {
	// This function is complete
	// fullPath := fmt.Sprintf("%s%s", patch, name)
	sFullPath := filepath.Join(sPath, sName)
	fDir, err := os.Stat(sFullPath)
	if err != nil {
		return false
	}

	return fDir.IsDir()
}

func Uint64ToBinary(u uint64) []byte {
	slB := make([]byte, 8)
	binary.BigEndian.PutUint64(slB, u)
	return slB
}

func BinaryToUint64(slB []byte) uint64 {
	u := binary.BigEndian.Uint64(slB)
	return u
}

func Encode64(sInput string) string { // sInput - ordinary string
	return base64.StdEncoding.EncodeToString([]byte(sInput))
}

func Decode64(sInput string) string { // sInput - base64 string
	slBData, err := base64.StdEncoding.DecodeString(sInput)
	if err != nil {
		return ""
	}
	return string(slBData)
}

func intPow(uAcc, uBase uint64, uExponent uint8) uint64 {
	// This function is complete
	if uExponent == 1 {
		return uAcc
	}
	return intPow(uAcc*uBase, uBase, uExponent-1)
}

func Pow(uBase uint64, uExponent uint8) uint64 {
	// This function is complete
	if uExponent == 0 {
		return 1
	}
	if uExponent == 1 {
		return uBase
	}
	return intPow(uBase*uBase, uBase, uExponent-1)
}
