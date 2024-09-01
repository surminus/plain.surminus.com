package main

import (
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
)

var serve bool
var port int

func init() {
	flag.BoolVar(&serve, "serve", false, "Serve the site")
	flag.IntVar(&port, "port", 80, "Port to serve on")
	flag.Parse()
}

const BuildDirectory = "docs"
const ContentDirectory = "content"

func main() {
	if err := os.RemoveAll(BuildDirectory); err != nil {
		log.Fatal(err)
	}
	if err := os.Mkdir(BuildDirectory, 0755); err != nil {
		log.Fatal(err)
	}

	var stylesheet string

	// If styles.css exists, copy it into the build directory and ensure it's
	// included in the generation
	if f, err := os.Stat("styles.css"); err == nil {
		data, err := os.ReadFile(f.Name())
		if err != nil {
			log.Fatal(err)
		}

		if err := os.WriteFile(filepath.Join(BuildDirectory, "styles.css"), data, 0644); err != nil {
			log.Fatal(err)
		}

		stylesheet = "/styles.css"
	}

	if err := filepath.WalkDir(ContentDirectory, func(path string, d fs.DirEntry, err error) error {
		if !strings.HasSuffix(path, ".md") {
			return nil
		}

		name := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
		destdir := filepath.Join(BuildDirectory, strings.TrimPrefix(filepath.Dir(path), ContentDirectory))

		if err := os.MkdirAll(destdir, 0755); err != nil {
			return err
		}

		filename := fmt.Sprintf("%s.html", name)
		dest := filepath.Join(destdir, filename)

		log.Println(dest)
		return NewPandoc(path, dest, stylesheet).Write()
	}); err != nil {
		log.Fatal(err)
	}

	if serve {
		log.Println("Serving at:", fmt.Sprintf("http://localhost:%d", port))
		cmd := exec.Command(docker(), []string{"run", "--rm", "--name", "surminus.com", "--volume", "./build/:/usr/share/nginx/html:ro", "--publish", fmt.Sprintf("%d:80", port), "nginx"}...)
		if err := cmd.Run(); err != nil {
			log.Fatal(err)
		}
	}
}

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

func docker() string {
	bin, err := exec.LookPath("docker")
	if err != nil {
		log.Fatal("docker missing")
	}

	return bin
}
