package vqlanalyzer

import (
	"errors"

	"github.com/Kwynto/GracefulDB/internal/engine/basicsystem/gauth"
	"github.com/Kwynto/GracefulDB/internal/engine/basicsystem/vqlexp"
)

func trimQuotationMarks(input string) string {
	if vqlexp.MRegExpCollection["QuotationMarks"].MatchString(input) {
		input = vqlexp.MRegExpCollection["QuotationMarks"].ReplaceAllLiteralString(input, "")
		return input
	}

	if vqlexp.MRegExpCollection["SpecQuotationMark"].MatchString(input) {
		input = vqlexp.MRegExpCollection["SpecQuotationMark"].ReplaceAllLiteralString(input, "")
		return input
	}

	return input
}

func preCheckerVQL(sTicket string) (sLogin string, stAccess gauth.TProfile, sNewTicket string, err error) {
	if sTicket == "" {
		return sLogin, stAccess, sNewTicket, errors.New("an empty ticket")
	}

	sLogin, stAccess, sNewTicket, err = gauth.CheckTicket(sTicket)
	if err != nil {
		return sLogin, stAccess, sNewTicket, err
	}

	if stAccess.Status.IsBad() {
		return sLogin, stAccess, sNewTicket, errors.New("auth error")
	}

	return sLogin, stAccess, sNewTicket, nil
}
