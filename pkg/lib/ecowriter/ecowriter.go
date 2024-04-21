package ecowriter

import (
	"bytes"
	"encoding/json"
	"os"
	"strings"
	"sync"
)

var (
	fEncoder *json.Encoder
	fDecoder *json.Decoder
	fBlock   sync.RWMutex
)

// Saving the structure to a JSON-file.
func WriteJSON(name string, data any) (err error) {
	// This function is complete
	fBlock.Lock()
	defer fBlock.Unlock()

	if _, err := os.Stat(name); !os.IsNotExist(err) {
		if err2 := os.Remove(name); err2 != nil {
			return err2
		}
	}

	wFile, err := os.OpenFile(name, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	defer wFile.Close()

	fEncoder = json.NewEncoder(wFile)
	if err := fEncoder.Encode(data); err != nil {
		return err
	}

	return nil
}

// Loading a structure from a JSON-file.
func ReadJSON(name string, data any) (err error) {
	// This function is complete
	fBlock.RLock()
	defer fBlock.RUnlock()

	if _, err := os.Stat(name); os.IsNotExist(err) {
		return err
	}

	rFile, err := os.OpenFile(name, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	defer rFile.Close()

	fDecoder = json.NewDecoder(rFile)
	if err := fDecoder.Decode(data); err != nil {
		return err
	}

	return nil
}

// Packing data in JSON-string
func EncodeJSON(data any) string {
	var buf bytes.Buffer

	encoder := json.NewEncoder(&buf)
	if err := encoder.Encode(data); err != nil {
		return ""
	}

	return buf.String()
}

// Getting data from JSON-string
func DecodeJSON(str string) any {
	var data any
	reader := strings.NewReader(str)

	decoder := json.NewDecoder(reader)
	if err := decoder.Decode(&data); err != nil {
		return nil
	}

	return data
}

func FileRead(name string) (string, error) {
	// This function is complete

	// f, err := os.Open(name)
	// if err != nil {
	// 	return "", err
	// }
	// defer f.Close()

	// buf := bytes.Buffer{}
	// sc := bufio.NewScanner(f)
	// for sc.Scan() {
	// 	buf.WriteString(sc.Text())
	// }

	// return buf.String(), nil

	bRead, err := os.ReadFile(name)
	if err != nil {
		return "", err
	}
	return string(bRead), nil
}
