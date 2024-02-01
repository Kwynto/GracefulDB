package sqlanalyzer

import "github.com/Kwynto/GracefulDB/pkg/lib/e"

// DDL — язык определения данных (Data Definition Language)

func (q *tQuery) DDLSearchCreate() (err error) {
	// -
	op := "internal -> analyzers -> sql -> DDL -> DDLSearchCreate"
	defer func() { e.Wrapper(op, err) }()

	return nil
}

func (q *tQuery) DDLSearchAlter() (err error) {
	// -
	op := "internal -> analyzers -> sql -> DDL -> DDLSearchAlter"
	defer func() { e.Wrapper(op, err) }()

	return nil
}

func (q *tQuery) DDLSearchDrop() (err error) {
	// -
	op := "internal -> analyzers -> sql -> DDL -> DDLSearchDrop"
	defer func() { e.Wrapper(op, err) }()

	return nil
}
