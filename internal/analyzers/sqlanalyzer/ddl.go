package sqlanalyzer

import "github.com/Kwynto/GracefulDB/pkg/lib/e"

// DDL — язык определения данных (Data Definition Language)

func (q *tQuery) DDLCreate() (err error) {
	// -
	op := "internal -> analyzers -> sql -> DDL -> DDLCreate"
	defer func() { e.Wrapper(op, err) }()

	return nil
}

func (q *tQuery) DDLAlter() (err error) {
	// -
	op := "internal -> analyzers -> sql -> DDL -> DDLAlter"
	defer func() { e.Wrapper(op, err) }()

	return nil
}

func (q *tQuery) DDLDrop() (err error) {
	// -
	op := "internal -> analyzers -> sql -> DDL -> DDLDrop"
	defer func() { e.Wrapper(op, err) }()

	return nil
}
