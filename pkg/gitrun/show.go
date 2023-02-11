package gitrun

func (repo *Repo) Show(file string) ([]byte, error) {
	// opts := []string{}
	// if s.repo != "" {
	// 	opts = append(opts, "-C", s.repo)
	// }
	// bpath := s.branch + "~" + strconv.Itoa(dp) + ":" + s.file
	// opts = append(opts, "show", bpath)

	// cmd := exec.Command("git", opts...)
	// out, err := cmd.Output()

	// out, err := Run("show", opts...)
	// if err != nil {
	// 	return nil, err
	// }

	out, err := repo.run("show", file)
	if err != nil {
		return nil, err
	}

	return out, nil
}
