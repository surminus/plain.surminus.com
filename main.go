package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/user"
)

func main() {
	NewPandoc("index.md", "index.html").Write()
}

type Pandoc struct {
	Source      string
	Destination string
}

func NewPandoc(src, dest string) *Pandoc {
	return &Pandoc{src, dest}
}

func (p *Pandoc) Write() error {
	cmd := p.DockerCmd()
	return cmd.Run()
}

func (p *Pandoc) DockerCmd(args ...string) *exec.Cmd {
	bin, err := exec.LookPath("docker")
	if err != nil {
		log.Fatal("docker missing")
	}

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	uid, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	cmdargs := []string{
		"run",
		"--rm",
		"--volume",
		fmt.Sprintf("%s:/data", cwd),
		"--user",
		fmt.Sprintf("%s:%s", uid.Uid, uid.Gid),
		"pandoc/latex",
		p.Source,
		"-s",
		"-o",
		p.Destination,
	}

	cmd := exec.Command(bin, append(cmdargs, args...)...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	return cmd
}
