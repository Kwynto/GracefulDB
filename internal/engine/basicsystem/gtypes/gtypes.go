package gtypes

import "time"

type Secret struct {
	Ticket   string `json:"ticket,omitempty"`
	Login    string `json:"login,omitempty"`
	Password string `json:"password,omitempty"`
	Hash     string `json:"hash,omitempty"`
}

type Response struct {
	State  string `json:"state,omitempty"`
	Ticket string `json:"ticket,omitempty"`
	Result string `json:"result,omitempty"`
}

type ResponseStrings struct {
	State  string   `json:"state,omitempty"`
	Ticket string   `json:"ticket,omitempty"`
	Result []string `json:"result,omitempty"`
}

type ResponseUints struct {
	State  string   `json:"state,omitempty"`
	Ticket string   `json:"ticket,omitempty"`
	Result []uint64 `json:"result,omitempty"`
}

type ResultColumn struct {
	Field      string    `json:"field"`
	Default    string    `json:"default"`
	NotNull    bool      `json:"notnull"`
	Unique     bool      `json:"unique"`
	LastUpdate time.Time `json:"lastupdate"`
}

type ResponseColumns struct {
	State  string         `json:"state,omitempty"`
	Ticket string         `json:"ticket,omitempty"`
	Result []ResultColumn `json:"result,omitempty"`
}

type TAccessFlags struct {
	Create bool `json:"create,omitempty"`
	Alter  bool `json:"alter,omitempty"`
	Drop   bool `json:"drop,omitempty"`
	Select bool `json:"select,omitempty"`
	Insert bool `json:"insert,omitempty"`
	Update bool `json:"update,omitempty"`
	Delete bool `json:"delete,omitempty"`
}

func (a TAccessFlags) AnyTrue() bool {
	return a.Create || a.Alter || a.Drop || a.Select || a.Insert || a.Update || a.Delete
}

type TAccess struct {
	Owner string                  `json:"owner,omitempty"` // login
	Flags map[string]TAccessFlags `json:"flags,omitempty"` // login - TAccessFlags
}

func DefaultSecret() Secret {
	return Secret{}
}
