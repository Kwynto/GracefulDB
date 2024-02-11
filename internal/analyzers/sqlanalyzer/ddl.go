package sqlanalyzer

import "github.com/Kwynto/GracefulDB/pkg/lib/e"

// DDL — язык определения данных (Data Definition Language)

func (q tQuery) DDLCreate() (result string, err error) {
	// -
	op := "internal -> analyzers -> sql -> DDL -> DDLCreate"
	defer func() { e.Wrapper(op, err) }()

	return "DDLCreate", nil
}

func (q tQuery) DDLAlter() (result string, err error) {
	// -
	op := "internal -> analyzers -> sql -> DDL -> DDLAlter"
	defer func() { e.Wrapper(op, err) }()

	return "DDLAlter", nil
}

func (q tQuery) DDLDrop() (result string, err error) {
	// -
	op := "internal -> analyzers -> sql -> DDL -> DDLDrop"
	defer func() { e.Wrapper(op, err) }()

	return "DDLDrop", nil
}
