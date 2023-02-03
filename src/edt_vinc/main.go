package main

import (
	"encoding/xml"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

// -r e:\dev\bhp_test -b dev -f src/Configuration/Configuration.mdo -d 10 -v
// -r e:\dev\bhp_test -b dev -f pl.txt -d 10 -v

type settings struct {
	repo    string
	branch  string
	file    string
	depth   int
	verbose bool
	change  bool
}

type Configuation struct {
	XMLName xml.Name `xml:"Configuration"`
	Version string   `xml:"version"`
}

var verbose bool

func main() {
	s := &settings{}
	s.parse()

	verbose = s.verbose

	err := run(s)
	if err != nil {
		log.Fatal(err)
	}

}

func run(s *settings) error {
	max, err := getMaxFromGit(s)
	if verbose {
		fmt.Println("Maximum version", max)
	}

	next, err := nextver(max)
	if err != nil {
		return err
	}
	if verbose {
		fmt.Println("Next version", next)
	}

	if s.change {
		changeVersion(s, next)
	}

	return nil
}

func (s *settings) parse() {
	flag.StringVar(&s.repo, "r", "", "target repository path")
	flag.StringVar(&s.branch, "b", "", "target branch")
	flag.StringVar(&s.file, "f", "", "target file")
	flag.IntVar(&s.depth, "d", 10, "search depth")
	flag.BoolVar(&s.verbose, "v", false, "verbose")
	flag.BoolVar(&s.change, "c", false, "change version in file")
	flag.Parse()
}

func getMaxFromGit(s *settings) (string, error) {
	max := "0.0.0.0"

	for i := 0; i < s.depth; i++ {
		c, err := gitShow(s, i)
		if err != nil {
			break
		}

		ver := versionFromConfig(c)

		if verbose {
			fmt.Print(s.branch, "~", i, " Compare ", max, " and ", ver, " -> ")
		}
		mver, err := maxver(max, ver)
		if err != nil {
			if verbose {
				fmt.Println(err.Error())
			}
			continue
		}

		max = mver
		if verbose {
			fmt.Println(mver)
		}
	}

	return max, nil
}

func changeVersion(s *settings, ver string) {
	p := "."
	if s.repo != "" {
		p = s.repo
	}
	fpath, _ := filepath.Abs(filepath.Join(p, s.file))

	if verbose {
		fmt.Println("Changing version in", fpath, "to", ver)
	}

	writeConfigurationVersion(fpath, ver)

}

func writeConfigurationVersion(filepath string, version string) {
	input, err := ioutil.ReadFile(filepath)
	if err != nil {
		log.Fatalln(err)
	}

	lines := strings.Split(string(input), "\n")

	for i, line := range lines {
		idx := strings.Index(line, "<version>")
		if idx > -1 {
			newline := line[:idx] + "<version>" + version + "</version>"
			lines[i] = newline
		}
	}
	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile(filepath, []byte(output), 0644)
	if err != nil {
		log.Fatalln(err)
	}
}

func maxver(ver1 string, ver2 string) (string, error) {
	ver1a, err := numver(ver1)
	if err != nil {
		return "", errors.New("cant' convert ver1 " + ver1)
	}

	ver2a, err := numver(ver2)
	if err != nil {
		return "", errors.New("cant' convert ver2 " + ver2)
	}

	for i := 0; i < 4; i++ {
		if ver2a[i] > ver1a[i] {
			return ver2, nil
		}
	}

	return ver1, nil
}

func nextver(ver string) (string, error) {
	vera, err := numver(ver)
	if err != nil {
		return "", errors.New("cant' convert ver " + ver)
	}

	var vers string
	for i := 0; i < 4; i++ {
		cnum := vera[i]
		if i == 3 {
			cnum++
		}
		vers = vers + strconv.Itoa(cnum)
		if i < 3 {
			vers = vers + "."
		}
	}

	return vers, nil
}

func numver(ver string) ([4]int, error) {
	var nver [4]int

	aver := strings.Split(ver, ".")
	if len(aver) != 4 {
		return nver, errors.New("wrong format")
	}

	for i := 0; i < 4; i++ {
		num, err := strconv.Atoi(aver[i])
		if err != nil {
			return nver, err
		}
		nver[i] = num
	}

	return nver, nil
}

func gitShow(s *settings, dp int) ([]byte, error) {
	opts := []string{}
	if s.repo != "" {
		opts = append(opts, "-C", s.repo)
	}
	bpath := s.branch + "~" + strconv.Itoa(dp) + ":" + s.file
	opts = append(opts, "show", bpath)

	cmd := exec.Command("git", opts...)
	out, err := cmd.Output()

	if err != nil {
		return nil, err
	}

	return out, nil
}

func versionFromConfig(c []byte) string {
	var conf Configuation
	var version string = "0.0.0.0"

	xml.Unmarshal(c, &conf)

	if conf.Version != "" {
		version = conf.Version
	}

	return version
}
