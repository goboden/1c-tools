package v8repo

import (
	"errors"
	"fmt"
	"os"
)

type Repository struct {
	Directory string
	Database  *Database
	Users     []User
}

type User struct {
	Name string
	GUID string
}

func NewRepository(dir string) (*Repository, error) {
	if info, err := os.Stat(dir); errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("directory '%s' not exist", dir)
	} else {
		if !info.IsDir() {
			return nil, fmt.Errorf("'%s' is not a directory", dir)
		}
	}

	rep := &Repository{Directory: dir}

	var err error
	rep.Database, err = NewDatabase(rep.Directory + "/1cv8ddb.1cd")
	if err != nil {
		return nil, err
	}

	// rep.ReadUsers()

	return rep, nil
}

func (rep *Repository) ReadUsers() {
	rep.Database.NewTableReader("USERS")
}
