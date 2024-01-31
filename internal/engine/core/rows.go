package core

type tCell struct {
	Name      string
	Value     []byte
	Operation rune
}

type tRow []tCell

func CreateRow(row tRow) uint64 {
	return 0
}

func ReadRow(id uint64) tRow {
	b := Uint64ToBinary(id)
	cell := tCell{
		Name:      POSTFIX_ID,
		Value:     b,
		Operation: '=',
	}

	return ReadRows(tRow{cell})
}

func ReadRows(filter tRow) tRow {
	return tRow{}
}

func UpdateRow(id uint64, row tRow) bool {
	b := Uint64ToBinary(id)
	cell := tCell{
		Name:      POSTFIX_ID,
		Value:     b,
		Operation: '=',
	}

	return UpdateRows(tRow{cell}, row)
}

func UpdateRows(filter tRow, row tRow) bool {
	return false
}

func DelRow(id uint64) bool {
	b := Uint64ToBinary(id)
	cell := tCell{
		Name:      POSTFIX_ID,
		Value:     b,
		Operation: '=',
	}

	return DelRows(tRow{cell})
}

func DelRows(filter tRow) bool {
	return false
}
