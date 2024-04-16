package core

import (
	"bufio"
	"bytes"
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

func FileRead(name string) (string, error) {
	// This function is complete

	// bRead, err := os.ReadFile(name)
	// if err != nil {
	// 	return "", err
	// }
	// return string(bRead), nil

	f, err := os.Open(name)
	if err != nil {
		return "", err
	}
	defer f.Close()

	buf := bytes.Buffer{}
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		buf.WriteString(sc.Text())
	}

	return buf.String(), nil
}
