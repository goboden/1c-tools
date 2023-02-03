package gitrun

import "os/exec"

type Repo struct {
	path string
}

func NewRepo(path string) *Repo {
	return &Repo{path}
}

func (repo *Repo) run(cmd string, opts ...string) ([]byte, error) {
	gitopts := []string{}

	gitopts = append(gitopts, cmd)

	if repo.path != "" {
		gitopts = append(gitopts, "-C", repo.path)
	}

	gitopts = append(gitopts, opts...)

	com := exec.Command("git", gitopts...)
	out, err := com.Output()

	if err != nil {
		return nil, err
	}

	return out, nil
}
