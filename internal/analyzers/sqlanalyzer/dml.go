package sqlanalyzer

import "github.com/Kwynto/GracefulDB/pkg/lib/e"

// DML — язык изменения данных (Data Manipulation Language)

func (q *tQuery) DMLSelect() (result string, err error) {
	// -
	op := "internal -> analyzers -> sql -> DML -> DMLSelect"
	defer func() { e.Wrapper(op, err) }()

	return "DMLSelect", nil
}

func (q *tQuery) DMLInsert() (result string, err error) {
	// -
	op := "internal -> analyzers -> sql -> DML -> DMLInsert"
	defer func() { e.Wrapper(op, err) }()

	return "DMLInsert", nil
}

func (q *tQuery) DMLUpdate() (result string, err error) {
	// -
	op := "internal -> analyzers -> sql -> DML -> DMLUpdate"
	defer func() { e.Wrapper(op, err) }()

	return "DMLUpdate", nil
}

func (q *tQuery) DMLDelete() (result string, err error) {
	// -
	op := "internal -> analyzers -> sql -> DML -> DMLDelete"
	defer func() { e.Wrapper(op, err) }()

	return "DMLDelete", nil
}

func (q *tQuery) DMLTruncate() (result string, err error) {
	// -
	op := "internal -> analyzers -> sql -> DML -> DMLTruncate"
	defer func() { e.Wrapper(op, err) }()

	return "DMLTruncate", nil
}

func (q *tQuery) DMLCommit() (result string, err error) {
	// -
	op := "internal -> analyzers -> sql -> DML -> DMLCommit"
	defer func() { e.Wrapper(op, err) }()

	return "DMLCommit", nil
}

func (q *tQuery) DMLRollback() (result string, err error) {
	// -
	op := "internal -> analyzers -> sql -> DML -> DMLRollback"
	defer func() { e.Wrapper(op, err) }()

	return "DMLRollback", nil
}
