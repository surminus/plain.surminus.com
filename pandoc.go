package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/user"
)

type Pandoc struct {
	Source      string
	Destination string
	Stylesheet  string
}

func NewPandoc(src, dest, stylesheet string) *Pandoc {
	return &Pandoc{src, dest, stylesheet}
}

func (p *Pandoc) Write() error {
	cmd := p.DockerCmd()
	return cmd.Run()
}

func (p *Pandoc) DockerCmd(args ...string) *exec.Cmd {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	uid, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	coreargs := []string{
		"run",
		"--rm",
		"--volume",
		fmt.Sprintf("%s:/data", cwd),
		"--user",
		fmt.Sprintf("%s:%s", uid.Uid, uid.Gid),
		"pandoc/latex",
		p.Source,
		"-f",
		"markdown+smart",
		"--to",
		"html5",
	}

	if p.Stylesheet != "" {
		coreargs = append(coreargs, "--css", p.Stylesheet)
	}

	cmdargs := append(coreargs, []string{"-s", "-o", p.Destination}...)

	cmd := exec.Command(docker(), append(cmdargs, args...)...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	return cmd
}
