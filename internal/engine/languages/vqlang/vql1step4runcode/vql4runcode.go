package vql1step4runcode

// Chain tabltes of simbols

type TVariableData struct {
	Type  int // 0 - nondetected (any), 1 - bool, 2 - char, 3 - byte, 4 - int, 5 - float, 6 - str, 7 - array, 8 - object, 9 - corteg
	Value string
}

type TPointer struct {
	Simbol string
	Link   *TTableOfSimbols
}

type TFunction struct {
	// Name     string
	Input    []TArgument
	FuncCode TActions
}

type TMapVariables map[string]TVariableData
type TMapPointers map[string]TPointer
type TMapFunctions map[string]TFunction

type TTableOfSimbols struct {
	Parent *TTableOfSimbols // only child table

	Variables TMapVariables
	Pointers  TMapPointers  // only child table
	Function  TMapFunctions // only root table

	// Input  TMapVariables // only root table
	Transparent bool // видимость родительской таблицы символов без указателей
}

// Production

type TReturn struct {
	Result   bool
	Returned bool
	Logic    bool
}

type TArgument struct {
	Productions TActions
	Simbol      string // имя в таблице символов
	Value       string // для переменной это значение, а для указателя это название переменной в другой ТС
	Link        *TTableOfSimbols
	Term        int // 0 - non, 1 - prodaction, 2 - variable, 3 - pointer
	Type        int // 0 - nondetected (any), 1 - bool, 2 - char, 3 - byte, 4 - int, 5 - float, 6 - str, 7 - array, 8 - object, 9 - corteg
}

type TProduction struct {
	// тип простого действия
	// - 0: системные базовые операции
	// -- 1: НЕкод или комментарий
	// -- 2: инструкция компилятору
	// -- 3: НЕтерминал
	// -- 4: терминал
	// --
	// -- 11: простая самостоятельная продукция
	// -- 12: блок области видимости
	// -- 13: функция
	// -- 14: возвратная операция
	// -- 15: условная операция
	// -- 16: классический цикл
	// -- 17: цикл по диапазону
	// --
	// -- 41: присваивание :=
	// -- 42: разделение голов и хвостов |=
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
	// название функций или директивы компилятору
	Name string
	// для кода блоков видимости и блочных структур
	LocalCode []TActions
	// наборр аргументов по-порядку из таблицы символов для действий и функций
	Left  []TArgument
	Right []TArgument
}

type TActions []TProduction

type TCode []string

type TComplexCode struct {
	Actions TActions
	Code    TCode
}
