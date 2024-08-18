package vqlex

import (
	"errors"

	"github.com/Kwynto/GracefulDB/internal/engine/basicsystem/gauth"
)

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
