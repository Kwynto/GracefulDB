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

func DefaultSecret() Secret {
	return Secret{}
}
