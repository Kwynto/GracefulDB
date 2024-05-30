package gtypes

import (
	"sync"
	"time"
)

type TColumnSpecification struct {
	Default string `json:"default"`
	NotNull bool   `json:"notnull"`
	Unique  bool   `json:"unique"` // FIXME: not used
}

type TColumnForWrite struct {
	Name    string
	OldName string
	Spec    TColumnSpecification

	// Flags of changes
	IsChName bool
}

type TColumnForStore struct {
	Field string
	Value string
}

type TRowForStore struct {
	Row    []TColumnForStore
	Id     uint64
	Time   int64
	Status int64 // memoried = 0  -  saved = 1  -  stored = 2
	Shape  int64 // primary = 0  -  required = 10  -  updated = 20  -  deleted = 30
	DB     string
	Table  string
}

type TWriteBuffer struct {
	Area     []TRowForStore
	BlockBuf sync.RWMutex
}

type TCollectBuffers struct {
	FirstBox  TWriteBuffer
	SecondBox TWriteBuffer
	Block     sync.RWMutex
	Switch    uint8
}

type TResponse struct {
	State  string `json:"state,omitempty"`
	Ticket string `json:"ticket,omitempty"`
	Result string `json:"result,omitempty"`
}

type TResponseStrings struct {
	Result []string `json:"result,omitempty"`
	State  string   `json:"state,omitempty"`
	Ticket string   `json:"ticket,omitempty"`
}

type TResponseUints struct {
	Result []uint64 `json:"result,omitempty"`
	State  string   `json:"state,omitempty"`
	Ticket string   `json:"ticket,omitempty"`
}

type TResponseRow map[string]string // name column and value

type TResponseSelect struct {
	Result []TResponseRow `json:"result,omitempty"`
	State  string         `json:"state,omitempty"`
	Ticket string         `json:"ticket,omitempty"`
}

type TResultColumn struct {
	Field      string    `json:"field"`
	Default    string    `json:"default"`
	NotNull    bool      `json:"notnull"`
	Unique     bool      `json:"unique"`
	LastUpdate time.Time `json:"lastupdate"`
}

type TResponseColumns struct {
	State  string          `json:"state,omitempty"`
	Ticket string          `json:"ticket,omitempty"`
	Result []TResultColumn `json:"result,omitempty"`
}

type TConditions struct {
	Type      string // "operation", "or", "and"
	Key       string
	Operation string
	Value     string
}

type TOrderBy struct {
	Cols []string
	Sort []uint8 // 0 - undef, 1 - asc, 2 - desc
}

type TUpdaateStruct struct {
	Where   []TConditions
	Couples map[string]string
}

type TSelectStruct struct {
	Orderby  TOrderBy
	Groupby  []string
	Where    []TConditions
	Columns  []string
	IsOrder  bool
	IsGroup  bool
	IsWhere  bool
	Distinct bool
}

type TDeleteStruct struct {
	Where   []TConditions
	IsWhere bool
}

type TAdditionalData struct {
	Db    string
	Table string
	Stamp int64 // for caching
}

type TSecret struct {
	Ticket   string `json:"ticket,omitempty"`
	Login    string `json:"login,omitempty"`
	Password string `json:"password,omitempty"`
	Hash     string `json:"hash,omitempty"`
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

func DefaultSecret() TSecret {
	return TSecret{}
}
