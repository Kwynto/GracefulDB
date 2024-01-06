package gtypes

type RowDB map[string]any

type VSecret struct {
	Ticket   string `json:"ticket,omitempty"`
	Login    string `json:"login,omitempty"`
	Password string `json:"password,omitempty"`
	Hash     string `json:"hash,omitempty"`
	QueryID  string `json:"queryid,omitempty"`
	Version  string `json:"version,omitempty"`
}

type VData []RowDB

type VFields RowDB

type VQuery struct {
	Action string  `json:"action"`
	Secret VSecret `json:"secret"`
	DB     string  `json:"db"`
	Table  string  `json:"table"`
	Fields VFields `json:"fields"`
	Data   VData   `json:"data"`
}

type VAnswer struct {
	Action      string  `json:"action"`
	Secret      VSecret `json:"secret,omitempty"`
	Data        VData   `json:"data,omitempty"`
	Error       int     `json:"error"`
	Description string  `json:"description,omitempty"`
}

func DefaultData() VData {
	return VData{}
}
