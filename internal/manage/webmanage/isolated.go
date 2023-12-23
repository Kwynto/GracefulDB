package webmanage

import (
	"net/http"
	"path/filepath"
)

// Isolation of statistical data for web access.

type isolatedFS struct {
	fs http.FileSystem
}

func (ifs isolatedFS) Open(path string) (http.File, error) {
	f, err := ifs.fs.Open(path)
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

// Isolation of authorization in handlers.

func IsolatedAuth() error {
	return nil
}
