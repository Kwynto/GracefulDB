package vqlex

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/Kwynto/GracefulDB/internal/engine/basicsystem/gauth"
	"github.com/Kwynto/GracefulDB/internal/engine/basicsystem/gtypes"
	"github.com/Kwynto/GracefulDB/internal/engine/languages/vqlang"
	"github.com/Kwynto/GracefulDB/pkg/lib/ecowriter"
)

type tInVariables struct {
	Name string
	Type string
}

type tStFuncCode struct {
	Name         string
	InVariables  []tInVariables
	OutVariables []string
	List         []string
}

type tQuery struct {
	Login          string
	Access         gauth.TProfile
	Ticket         string
	DB             string // DB name
	Table          string // Table name
	Code           vqlang.TCode
	Actions        vqlang.TActions
	LocalFunctions map[string]tStFuncCode
	TableOfSimbols vqlang.TTableOfSimbols
}

func splitCode(sOriginalCode string) vqlang.TCode {
	// This function is complete
	slList := strings.Split(sOriginalCode, "\n")

	return slList
}

func analyzer(slStCode vqlang.TCode) vqlang.TActions {
	// -
	_ = slStCode
	return vqlang.TActions{}
}

// TODO: Request
func Request(sTicket string, sOriginalCode string, sVariables string) string {
	// -
	var stRes gtypes.TResponse

	// Pre checking
	sLogin, stAccess, sNewTicket, errC := preCheckerVQL(sTicket)
	if errC != nil {
		stRes.State = "error"
		stRes.Result = errC.Error()
		slog.Debug("Wrong request:", slog.String("err", stRes.Result))
		return ecowriter.EncodeJSON(stRes)
	}

	if sNewTicket != "" {
		stRes.Ticket = sNewTicket
		sTicket = sNewTicket
	} else {
		stRes.Ticket = sTicket
	}

	// Table of simbols
	mVariables, okU := ecowriter.DecodeJSONMap(sVariables)
	if !okU {
		stRes.State = "error"
		stRes.Result = "invalid variable format"
		slog.Debug("Wrong request:", slog.String("err", stRes.Result))
		return ecowriter.EncodeJSON(stRes)
	}

	for sKey, inValue := range mVariables {
		rKey := []rune(sKey)[0]
		if rKey != rune('$') {
			newKey := fmt.Sprintf("$%s", sKey)
			mVariables[newKey] = inValue
			delete(mVariables, sKey)
		}
	}

	// Preparation query
	slStCode := splitCode(sOriginalCode)
	stActions := analyzer(slStCode)
	// FIXME: it
	// slQryLines, mLocalFunctions, errP := preparation(sOriginalCode)
	// if errP != nil {
	// 	stRes.State = "error"
	// 	stRes.Result = errP.Error()
	// 	slog.Debug("Wrong request:", slog.String("err", stRes.Result))
	// 	return ecowriter.EncodeJSON(stRes)
	// }

	var query tQuery = tQuery{
		Login:   sLogin,
		Access:  stAccess,
		Ticket:  sTicket,
		Code:    slStCode,
		Actions: stActions,
		// LocalFunctions: mLocalFunctions,
		TableOfSimbols: vqlang.TTableOfSimbols{
			Variables:   mVariables,
			Transparent: true,
		},
	}

	// FIXME: it
	_ = query

	return ecowriter.EncodeJSON(stRes)
}
