package sqlanalyzer

import "github.com/Kwynto/GracefulDB/pkg/lib/e"

// DML — язык изменения данных (Data Manipulation Language)

func (q *tQuery) DMLSelect() (result string, err error) {
	// -
	op := "internal -> analyzers -> sql -> DML -> DMLSearchSelect"
	defer func() { e.Wrapper(op, err) }()

	return "", nil
}

func (q *tQuery) DMLInsert() (result string, err error) {
	// -
	op := "internal -> analyzers -> sql -> DML -> DMLSearchInsert"
	defer func() { e.Wrapper(op, err) }()

	return "", nil
}

func (q *tQuery) DMLUpdate() (result string, err error) {
	// -
	op := "internal -> analyzers -> sql -> DML -> DMLSearchUpdate"
	defer func() { e.Wrapper(op, err) }()

	return "", nil
}

func (q *tQuery) DMLDelete() (result string, err error) {
	// -
	op := "internal -> analyzers -> sql -> DML -> DMLSearchDelete"
	defer func() { e.Wrapper(op, err) }()

	return "", nil
}

func (q *tQuery) DMLTruncate() (result string, err error) {
	// -
	op := "internal -> analyzers -> sql -> DML -> DMLSearchTruncate"
	defer func() { e.Wrapper(op, err) }()

	return "", nil
}

func (q *tQuery) DMLCommit() (result string, err error) {
	// -
	op := "internal -> analyzers -> sql -> DML -> DMLSearchCommit"
	defer func() { e.Wrapper(op, err) }()

	return "", nil
}

func (q *tQuery) DMLRollback() (result string, err error) {
	// -
	op := "internal -> analyzers -> sql -> DML -> DMLSearchRollback"
	defer func() { e.Wrapper(op, err) }()

	return "", nil
}
