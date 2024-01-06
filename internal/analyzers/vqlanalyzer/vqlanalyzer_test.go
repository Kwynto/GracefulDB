package vqlanalyzer

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/Kwynto/GracefulDB/internal/engine/basicsystem/gauth"
	"github.com/Kwynto/GracefulDB/internal/engine/basicsystem/gtypes"
)

func Test_Reading(t *testing.T) {
	auth := VAuth{
		Login:  "",
		Access: gauth.TProfile{},
	}
	fields := gtypes.VFields{}

	res := Reading(&auth, &fields)
	if reflect.TypeOf(res) != reflect.TypeOf(&gtypes.VData{}) {
		t.Error("Reading() error = The function returns the wrong type")
	}
}

func Test_Storing(t *testing.T) {
	auth := VAuth{
		Login:  "",
		Access: gauth.TProfile{},
	}
	in := gtypes.VQuery{}

	res := Storing(&auth, &in)
	if reflect.TypeOf(res) != reflect.TypeOf(&gtypes.VData{}) {
		t.Error("Storing() error = The function returns the wrong type")
	}
}

func Test_Deleting(t *testing.T) {
	auth := VAuth{
		Login:  "",
		Access: gauth.TProfile{},
	}
	fields := gtypes.VFields{}

	res := Deleting(&auth, &fields)
	if reflect.TypeOf(res) != reflect.TypeOf(&gtypes.VData{}) {
		t.Error("Deleting() error = The function returns the wrong type")
	}
}

func Test_Managing(t *testing.T) {
	auth := VAuth{
		Login:  "",
		Access: gauth.TProfile{},
	}
	in := gtypes.VQuery{}

	res := Managing(&auth, &in)
	if reflect.TypeOf(res) != reflect.TypeOf(&gtypes.VData{}) {
		t.Error("Managing() error = The function returns the wrong type")
	}
}

func Test_Request(t *testing.T) {
	instruction := "not valid instruction"
	in := []byte(instruction)
	res1 := Request(&in)
	if reflect.TypeOf(res1) != reflect.TypeOf(&[]byte{}) {
		t.Error("Request() error = The function returns the wrong type")
	}

	instruction = "{\"action\": true, \"db\": 1, \"secret\": \"secret\"}"
	in = []byte(instruction)
	res2 := Request(&in)
	if reflect.TypeOf(res2) != reflect.TypeOf(&[]byte{}) {
		t.Error("Request() error = The function returns the wrong type")
	}

	instruction = `{"action":"auth", "secret":{"login":"root", "password":"toor", "queryid":"any-id"}}`
	in = []byte(instruction)
	res3 := Request(&in)
	if reflect.TypeOf(res3) != reflect.TypeOf(&[]byte{}) {
		t.Error("Request() error = The function returns the wrong type")
	}
}

func Test_Processing(t *testing.T) {
	gauth.AuthFile = "./../../../config/auth.json"
	gauth.AccessFile = "./../../../config/access.json"
	gauth.Start()

	in := gtypes.VQuery{}
	res := Processing(&in)
	if reflect.TypeOf(res) != reflect.TypeOf(&gtypes.VAnswer{}) {
		t.Error("Processing() error = The function returns the wrong type")
	}

	var qry gtypes.VQuery

	instruction := `{"action":"auth", "secret":{"login":"root", "password":"toor", "queryid":"any-id"}}`
	inst := []byte(instruction)
	json.Unmarshal(inst, &qry)
	res1 := Processing(&qry)
	if reflect.TypeOf(res1) != reflect.TypeOf(&gtypes.VAnswer{}) {
		t.Error("Request() error = The function returns the wrong type")
	}

	// fmt.Println("Ticket:", res1.Secret.Ticket)

	instruction = `{"action":"read", "secret":{"ticket": "%s", "queryid":"any-id"}, "fields":{}}`
	instruction = fmt.Sprintf(instruction, res1.Secret.Ticket)
	inst = []byte(instruction)
	json.Unmarshal(inst, &qry)
	res2 := Processing(&qry)
	if reflect.TypeOf(res2) != reflect.TypeOf(&gtypes.VAnswer{}) {
		t.Error("Request() error = The function returns the wrong type")
	}

	instruction = `{"action":"store", "secret":{"ticket": "%s", "queryid":"any-id"}}`
	instruction = fmt.Sprintf(instruction, res1.Secret.Ticket)
	inst = []byte(instruction)
	json.Unmarshal(inst, &qry)
	res3 := Processing(&qry)
	if reflect.TypeOf(res3) != reflect.TypeOf(&gtypes.VAnswer{}) {
		t.Error("Request() error = The function returns the wrong type")
	}

	instruction = `{"action":"delete", "secret":{"ticket": "%s", "queryid":"any-id"}, "fields":{}}`
	instruction = fmt.Sprintf(instruction, res1.Secret.Ticket)
	inst = []byte(instruction)
	json.Unmarshal(inst, &qry)
	res4 := Processing(&qry)
	if reflect.TypeOf(res4) != reflect.TypeOf(&gtypes.VAnswer{}) {
		t.Error("Request() error = The function returns the wrong type")
	}

	instruction = `{"action":"manage", "secret":{"ticket": "%s", "queryid":"any-id"}}`
	instruction = fmt.Sprintf(instruction, res1.Secret.Ticket)
	inst = []byte(instruction)
	json.Unmarshal(inst, &qry)
	res5 := Processing(&qry)
	if reflect.TypeOf(res5) != reflect.TypeOf(&gtypes.VAnswer{}) {
		t.Error("Request() error = The function returns the wrong type")
	}

	instruction = `{"action":"", "secret":{"ticket": "%s", "queryid":"any-id"}}`
	instruction = fmt.Sprintf(instruction, res1.Secret.Ticket)
	inst = []byte(instruction)
	json.Unmarshal(inst, &qry)
	res6 := Processing(&qry)
	if reflect.TypeOf(res6) != reflect.TypeOf(&gtypes.VAnswer{}) {
		t.Error("Request() error = The function returns the wrong type")
	}

	instruction = `{"action":"fake-command", "secret":{"ticket": "%s", "queryid":"any-id"}}`
	instruction = fmt.Sprintf(instruction, res1.Secret.Ticket)
	inst = []byte(instruction)
	json.Unmarshal(inst, &qry)
	res7 := Processing(&qry)
	if reflect.TypeOf(res7) != reflect.TypeOf(&gtypes.VAnswer{}) {
		t.Error("Request() error = The function returns the wrong type")
	}

	instruction = `{"action":"auth", "secret":{"login":"root", "password":"toor", "queryid":"any-id"}}`
	inst = []byte(instruction)
	json.Unmarshal(inst, &qry)
	Processing(&qry)
	instruction = `{"action":"read", "secret":{"ticket": "%s", "queryid":"any-id"}, "fields":{}}`
	instruction = fmt.Sprintf(instruction, res1.Secret.Ticket)
	inst = []byte(instruction)
	json.Unmarshal(inst, &qry)
	res8 := Processing(&qry)
	if reflect.TypeOf(res8) != reflect.TypeOf(&gtypes.VAnswer{}) {
		t.Error("Request() error = The function returns the wrong type")
	}

}
