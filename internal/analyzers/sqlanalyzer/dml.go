package sqlanalyzer

import "github.com/Kwynto/GracefulDB/pkg/lib/e"

// DML — язык изменения данных (Data Manipulation Language)

func (q *tQuery) DMLSelect() (err error) {
	// -
	op := "internal -> analyzers -> sql -> DML -> DMLSelect"
	defer func() { e.Wrapper(op, err) }()

	return nil
}

func (q *tQuery) DMLInsert() (err error) {
	// -
	op := "internal -> analyzers -> sql -> DML -> DMLInsert"
	defer func() { e.Wrapper(op, err) }()

	return nil
}

func (q *tQuery) DMLUpdate() (err error) {
	// -
	op := "internal -> analyzers -> sql -> DML -> DMLUpdate"
	defer func() { e.Wrapper(op, err) }()

	return nil
}

func (q *tQuery) DMLDelete() (err error) {
	// -
	op := "internal -> analyzers -> sql -> DML -> DMLDelete"
	defer func() { e.Wrapper(op, err) }()

	return nil
}

func (q *tQuery) DMLTruncate() (err error) {
	// -
	op := "internal -> analyzers -> sql -> DML -> DMLTruncate"
	defer func() { e.Wrapper(op, err) }()

	return nil
}

func (q *tQuery) DMLCommit() (err error) {
	// -
	op := "internal -> analyzers -> sql -> DML -> DMLCommit"
	defer func() { e.Wrapper(op, err) }()

	return nil
}

func (q *tQuery) DMLRollback() (err error) {
	// -
	op := "internal -> analyzers -> sql -> DML -> DMLRollback"
	defer func() { e.Wrapper(op, err) }()

	return nil
}
