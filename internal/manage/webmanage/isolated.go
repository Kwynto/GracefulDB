package webmanage

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/Kwynto/GracefulDB/internal/engine/basicsystem/gauth"
	"github.com/Kwynto/GracefulDB/pkg/lib/e"
	"github.com/Kwynto/gosession"
)

// Isolation of statistical data for web access.

type isolatedFS struct {
	fs http.FileSystem
}

func (ifs isolatedFS) Open(path string) (f http.File, err error) {
	op := "internal -> WebManage -> isolated -> Open"
	defer func() { e.Wrapper(op, err) }()

	f, err = ifs.fs.Open(path)
	if err != nil {
		return nil, err
	}

	s, err := f.Stat()
	if err != nil {
		return nil, err
	}
	if s.IsDir() {
		index := filepath.Join(path, "index.html")
		if _, err := ifs.fs.Open(index); err != nil {
			closeErr := f.Close()
			if closeErr != nil {
				return nil, closeErr
			}
			return nil, err
		}
	}

	return f, nil
}

// Isolation of authorization.
func IsolatedAuth(w http.ResponseWriter, r *http.Request, rules []gauth.TRole) bool {
	sesID := gosession.Start(&w, r)
	auth := sesID.Get("auth")
	login := fmt.Sprint(auth)
	profile, err := gauth.GetProfile(login)
	if err != nil {
		return true
	}

	if !profile.IsAllowed(rules) {
		return true
	}

	if profile.Status == gauth.NEW {
		profile.Status = gauth.ACTIVE
		gauth.UpdateProfile(login, profile)
	}

	return false
}
