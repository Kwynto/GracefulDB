package gtypes

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

type TAccessFlags struct {
	Create bool `json:"create,omitempty"`
	Read   bool `json:"read,omitempty"`
	Update bool `json:"update,omitempty"`
	Delete bool `json:"delete,omitempty"`
}

type TAccess struct {
	Owner string                  `json:"owner,omitempty"` // login
	Flags map[string]TAccessFlags `json:"flags,omitempty"` // login - TAccessFlags
}

func DefaultSecret() Secret {
	return Secret{}
}
