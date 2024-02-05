package sqlanalyzer

import "github.com/Kwynto/GracefulDB/pkg/lib/e"

// DDL — язык определения данных (Data Definition Language)

func (q *tQuery) DDLCreate() (result string, err error) {
	// -
	op := "internal -> analyzers -> sql -> DDL -> DDLSearchCreate"
	defer func() { e.Wrapper(op, err) }()

	return "", nil
}

func (q *tQuery) DDLAlter() (result string, err error) {
	// -
	op := "internal -> analyzers -> sql -> DDL -> DDLSearchAlter"
	defer func() { e.Wrapper(op, err) }()

	return "", nil
}

func (q *tQuery) DDLDrop() (result string, err error) {
	// -
	op := "internal -> analyzers -> sql -> DDL -> DDLSearchDrop"
	defer func() { e.Wrapper(op, err) }()

	return "", nil
}
