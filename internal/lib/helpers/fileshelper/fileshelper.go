package fileshelper

import "os"

func FileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}

	return true
}

func CreateFile(name string) error {
	fo, err := os.Create(name)
	if err != nil {
		return err
	}

	fo.Close()

	return nil
}
