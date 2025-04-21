package gtypes

import (
	"sync"
	"time"
)

// Chain tabltes of simbols

// type TMapVariables map[string]any
type TVariableData struct {
	Type  int // 0 - nondetected (any), 1 - bool, 2 - char, 3 - byte, 4 - int, 5 - float, 6 - str, 7 - array, 8 - object, 9 - corteg
	Value string
}
type TPointer struct {
	Simbol string
	Link   *TTableOfSimbols
}
type TMapVariables map[string]TVariableData
type TMapPointers map[string]TPointer

type TTableOfSimbols struct {
	Parent *TTableOfSimbols // only child table

	Variables TMapVariables
	Pointers  TMapPointers // only child table

	Input  TMapVariables // only root table
	IsRoot bool
}

func (tos *TTableOfSimbols) Record(dest *TArgument) bool {
	if _, okV := tos.Variables[dest.Simbol]; okV {
		tos.Variables[dest.Simbol] = TVariableData{
			Type:  dest.Type,
			Value: dest.Value,
		}
	} else if pointer, okP := tos.Pointers[dest.Simbol]; okP {
		pointer.Link.Variables[pointer.Simbol] = TVariableData{
			Type:  dest.Type,
			Value: dest.Value,
		}
	} else if !tos.IsRoot {
		return tos.Parent.Record(dest)
	} else {
		return false
	}

	return true
}

func (tos *TTableOfSimbols) Set(dest *TArgument, source *TArgument) bool {
	switch source.Term {
	case 0:
		dest.Type = source.Type
		dest.Value = source.Value
	case 1:
		ok, resultArgument := source.Production.Exec(tos)
		if !ok {
			return false
		}
		dest.Value = resultArgument.Value
	case 2:
		dest.Term = source.Term
		dest.Type = source.Type
		dest.Value = source.Value
	case 3:
		dest.Term = source.Term
		dest.Value = source.Value
		dest.Link = source.Link
	}

	if dest.Simbol != "" {
		if !tos.Record(dest) {
			switch dest.Term {
			case 2:
				tos.Variables[dest.Simbol] = TVariableData{
					Type:  dest.Type,
					Value: dest.Value,
				}
			case 3:
				tos.Pointers[dest.Simbol] = TPointer{
					Simbol: dest.Value,
					Link:   dest.Link,
				}
			}
		}
	} else {
		return false
	}

	return true
}

// Lex Analyzer

type TArgument struct {
	Production TProduction
	// Name       string // для именованных аргументов
	Simbol string // имя в таблице символов
	Value  string // для переменной это значение, а для указателя это название переменной в другой ТС
	Link   *TTableOfSimbols
	Term   int // 0 - non, 1 - prodaction, 2 - variable, 3 - pointer
	Type   int // 0 - nondetected (any), 1 - bool, 2 - char, 3 - byte, 4 - int, 5 - float, 6 - str, 7 - array, 8 - object, 9 - corteg
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
	// -- 15: условная операция
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
	// -- 325: умножение с присвоением :*
	// -- 326: деление с присвоением :/
	// -- 327: остаток от деления с присвоением :%
	// -- 328: возведение в степень с присвоением :**
	// -- 329: извлечение корня с присвоением :/*
	// -- 330: логарифм с присвоением :/%
	// - 400: операции со строками
	// -- 401: конкатенация строк
	// - 500: (зарезервировано)
	// - 600: (дириктивы)
	// - 700: (зарезервировано)
	// - 800: (встроенные функции)
	Type int
	// наборр аргументов по-порядку из таблицы символов для действий и функций
	Dest   []TArgument
	Source []TArgument
}

type TActions []TProduction

func (actions TActions) Run(input TMapVariables) bool {
	pRootTOS := &TTableOfSimbols{
		Input:  input,
		IsRoot: true,
	}
	for _, production := range actions {
		if ok, _ := production.Exec(pRootTOS); !ok {
			return false
		}
	}
	return true
}

func (parentProduction TProduction) Exec(parentTOS *TTableOfSimbols) (bool, *TArgument) {
	// pLocalTOS := &TTableOfSimbols{
	// 	Parent: parentTOS,
	// 	IsRoot: false,
	// }
	switch parentProduction.Type {
	case 1, 2, 3, 4:
		// пока пропускаем
	case 11:
		zeroArgument := parentProduction.Dest[0]
		if zeroArgument.Term == 1 {
			if ok, _ := zeroArgument.Production.Exec(parentTOS); !ok {
				return false, &TArgument{}
			}
		}
	case 12, 13, 14, 15, 16, 17:
		// пока пропустить
	case 41:
		if len(parentProduction.Dest) == len(parentProduction.Source) {
			for indOA, outputArg := range parentProduction.Source {
				if !parentTOS.Set(&parentProduction.Dest[indOA], &outputArg) {
					return false, &TArgument{}
				}
			}
		} else {
			return false, &TArgument{}
		}
	}

	return true, &TArgument{}
}

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
