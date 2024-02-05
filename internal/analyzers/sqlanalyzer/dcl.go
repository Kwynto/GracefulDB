package sqlanalyzer

import (
	"github.com/Kwynto/GracefulDB/pkg/lib/e"
)

// DCL — язык управления данными (Data Control Language)

func (q *tQuery) DCLGrant() (result string, err error) {
	// -
	op := "internal -> analyzers -> sql -> DCL -> DCLSearchGrant"
	defer func() { e.Wrapper(op, err) }()

	return "", nil
}

func (q *tQuery) DCLRevoke() (result string, err error) {
	// -
	op := "internal -> analyzers -> sql -> DCL -> DCLSearchRevoke"
	defer func() { e.Wrapper(op, err) }()

	return "", nil
}

func (q *tQuery) DCLUse() (result string, err error) {
	// -
	op := "internal -> analyzers -> sql -> DCL -> DCLSearchUse"
	defer func() { e.Wrapper(op, err) }()

	return "", nil
}
