package gitrun

func (repo *Repo) Show(file string) ([]byte, error) {
	out, err := repo.run("show", file)
	if err != nil {
		return nil, err
	}

	return out, nil
}
