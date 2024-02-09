package sqlanalyzer

import (
	"github.com/Kwynto/GracefulDB/pkg/lib/e"
)

// DCL — язык управления данными (Data Control Language)

func (q *tQuery) DCLGrant() (result string, err error) {
	// -
	op := "internal -> analyzers -> sql -> DCL -> DCLGrant"
	defer func() { e.Wrapper(op, err) }()

	return "", nil
}

func (q *tQuery) DCLRevoke() (result string, err error) {
	// -
	op := "internal -> analyzers -> sql -> DCL -> DCLRevoke"
	defer func() { e.Wrapper(op, err) }()

	return "", nil
}

func (q *tQuery) DCLUse() (result string, err error) {
	// -
	op := "internal -> analyzers -> sql -> DCL -> DCLUse"
	defer func() { e.Wrapper(op, err) }()

	return "", nil
}

func (q *tQuery) DCLAuth() (result string, err error) {
	// -
	op := "internal -> analyzers -> sql -> DCL -> DCLAuth"
	defer func() { e.Wrapper(op, err) }()

	return "", nil
}
