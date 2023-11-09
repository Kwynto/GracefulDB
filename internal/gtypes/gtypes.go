package gtypes

type VQuery struct {
	Action string `json:"action"`
	Secret string `json:"secret"`
	DB     string `json:"db"`
	Table  string `json:"table"`
	Fields string `json:"fields"`
	Data   string `json:"data"`
}

type VAnswer struct {
	Action string `json:"action"`
	Secret string `json:"secret,omitempty"`
	Data   string `json:"data,omitempty"`
	Error  int    `json:"error"`
}
