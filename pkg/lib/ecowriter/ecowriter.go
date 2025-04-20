package ecowriter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/Kwynto/GracefulDB/internal/engine/basicsystem/gtypes"
)

var (
	jeEncoder *json.Encoder
	jdDecoder *json.Decoder
	mxBlock   sync.RWMutex
)

// Saving the structure to a JSON-file.
func WriteJSON(name string, data any) (err error) {
	// This function is complete
	mxBlock.Lock()
	defer mxBlock.Unlock()

	if _, errStat := os.Stat(name); !os.IsNotExist(errStat) {
		if errRemove := os.Remove(name); errRemove != nil {
			return errRemove
		}
	}

	f, err := os.OpenFile(name, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	defer f.Close()

	jeEncoder = json.NewEncoder(f)
	if err := jeEncoder.Encode(data); err != nil {
		return err
	}

	return nil
}

// Loading a structure from a JSON-file.
func ReadJSON(name string, data any) (err error) {
	// This function is complete
	mxBlock.RLock()
	defer mxBlock.RUnlock()

	if _, errStat := os.Stat(name); os.IsNotExist(errStat) {
		return errStat
	}

	f, err := os.OpenFile(name, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	defer f.Close()

	jdDecoder = json.NewDecoder(f)
	if err := jdDecoder.Decode(data); err != nil {
		return err
	}

	return nil
}

// Packing data in JSON-string
func EncodeJSON(inData any) string {
	var buf bytes.Buffer

	je := json.NewEncoder(&buf)
	if err := je.Encode(inData); err != nil {
		return ""
	}

	return buf.String()
}

// Getting data from JSON-string
func DecodeJSON(str string) any {
	var inData any
	reader := strings.NewReader(str)

	jd := json.NewDecoder(reader)
	if err := jd.Decode(&inData); err != nil {
		return nil
	}

	return inData
}

// Getting map from JSON-string
func DecodeJSONMap(str string) (gtypes.TMapVariables, bool) {
	var inData map[string]any
	// var outData map[string]string
	var outData gtypes.TMapVariables

	reader := strings.NewReader(str)

	jd := json.NewDecoder(reader)
	if err := jd.Decode(&inData); err != nil {
		return nil, false
	}

	for key, value := range inData {
		outData[key] = gtypes.TVariableData{
			// TODO: сделать проверку и приведение типов
			Type:  0,
			Value: fmt.Sprint(value),
		}

	}

	return outData, true
}

func FileRead(name string) (string, error) {
	bRead, err := os.ReadFile(name)
	if err != nil {
		return "", err
	}
	return string(bRead), nil
}
