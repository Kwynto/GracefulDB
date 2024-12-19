package gtypes

import (
	"sync"
	"time"
)

// Chain tabltes of simbols

type TMapVariables map[string]any
type TMapPointers map[string]*TTableOfSimbols

type TTableOfSimbols struct {
	Parent   *TTableOfSimbols // only child table
	Pointers TMapPointers     // only child table

	Input TMapVariables // only root table
	Self  TMapVariables

	IsRoot bool
}

// Lex Analyzer

type TArgument struct {
	Production TProduction
	Pointer    *TTableOfSimbols
	Name       string // для именованных аргументов
	Any        string
	Str        string
	Int        int
	Float      float64
	Type       int // 0 - non, 1 - prodaction, 2 - pointer, 12 - any, 11 - str, 12 - int, 13 - float
}

type TProduction struct {
	// тип простого действия
	// - 0: системные базовые операции
	// -- 1: НЕкод или комментарий
	// -- 2: инструкция компилятору
	// -- 3: НЕтерминал
	// -- 4: терминал
	// --
	// -- 11: простая вложенная продукция
	// -- 12: блок области видимости
	// -- 13: функция
	// -- 14: возвратная операция
	// -- 15: условная уперация
	// -- 16: классический цикл
	// -- 17: цикл по диапазону
	// --
	// -- 41: присваивание :=
	// --
	// - 100: логические операции
	// -- 101: равенство ==
	// -- 102: НЕравенство !=
	// -- 103: меньше <
	// -- 104: больше >
	// -- 105: меньше или равно <=
	// -- 106: больше или равно >=
	// --
	// -- 111: NOT
	// -- 112: AND
	// -- 113: OR
	// -- 114: XOR
	// --
	// - 200: битовые операции
	// -- 201: сдвиг влево с кольцевой заменой <<
	// -- 202: сдвиг вправо с кольцевой заменой >>
	// -- 203: сдвиг влево с заполнением нулем  <:<
	// -- 204: сдвиг вправо с заполнением нулем >:>
	// -- 205: сдвиг влево с заполнением единицей  <!<
	// -- 206: сдвиг вправо с заполнением единицей >!>
	// --
	// -- 211: кольцевой сдвиг влево с прсвоением :<<
	// -- 212: кольцевой сдвиг вправо с присвоением :>>
	// -- 213: сдвиг влево с заполнением нулем и с присвоением :<:<
	// -- 204: сдвиг вправо с заполнением нулем и с присвоением :>:>
	// -- 205: сдвиг влево с заполнением единицей и с присвоением :<!<
	// -- 206: сдвиг вправо с заполнением единицей и с присвоением :>!>
	// --
	// - 300: арифметические операции
	// -- 301: сложение +
	// -- 302: вычитание -
	// -- 303: умножение *
	// -- 304: деление /
	// -- 305: остаток от деления %
	// -- 306: возведение в степень **
	// -- 307: извлечение корня /*
	// -- 308: логарифм /%
	// --
	// -- 321: инкремет ++
	// -- 322: дикремент --
	// -- 323: инкремент с присвоением :++
	// -- 324: дикремент с присвоением :--
	// -- 325: умножение с присвоением :**
	// -- 326: деление с присвоением ://
	// -- 327: остаток от деления с присвоением :%%
	// -- 328: возведение в степень с присвоением :**
	// -- 329: извлечение корня с присвоением :/*
	// -- 330: логарифм с присвоением :/%
	// - 400: операции со строками
	// -- 401: конкатенация строк
	Type int
	// идентификатор для сложных продакций, например, номер вложенного блока для открывающих и закрывающих скобок
	// Id int
	// наборр аргументов по-порядку из таблицы символов для действий и функций
	Input  []TArgument
	Output []TArgument
}

type TActions []TProduction

type TCode []string

// Engine

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
	Is   bool
}

type TLimit struct {
	Start  int
	Offset int
	Is     bool
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
