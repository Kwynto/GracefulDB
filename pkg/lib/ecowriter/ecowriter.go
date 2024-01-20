package ecowriter

import (
	"encoding/json"
	"os"

	"github.com/Kwynto/GracefulDB/pkg/lib/e"
)

var (
	encoder *json.Encoder
	decoder *json.Decoder
)

func WriteJSON(name string, data any) (err error) {
	op := "pkg -> lib -> ecowriter -> WriteJSON"
	defer func() { e.Wrapper(op, err) }()

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

	encoder = json.NewEncoder(wFile)
	if err := encoder.Encode(data); err != nil {
		return err
	}

	return nil
}

func ReadJSON(name string, data any) (err error) {
	op := "pkg -> lib -> ecowriter -> ReadJSON"
	defer func() { e.Wrapper(op, err) }()

	if _, err := os.Stat(name); os.IsNotExist(err) {
		return err
	}

	rFile, err := os.OpenFile(name, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	defer rFile.Close()

	decoder = json.NewDecoder(rFile)
	if err := decoder.Decode(&data); err != nil {
		return err
	}

	return nil
}
