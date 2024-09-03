package main

import (
	"flag"
	"fmt"
	"log"
	"os/exec"
)

var serve bool
var port int

func init() {
	flag.BoolVar(&serve, "serve", false, "Serve the site")
	flag.IntVar(&port, "port", 80, "Port to serve on")
	flag.Parse()
}

const BuildDirectory = "build"
const ContentDirectory = "content"

func main() {
	if err := build(); err != nil {
		log.Fatal(err)
	}

	if serve {
		log.Println("Serving at:", fmt.Sprintf("http://localhost:%d", port))
		cmd := exec.Command(docker(), []string{"run", "--rm", "--name", "surminus.com", "--volume", fmt.Sprintf("./%s/:/usr/share/nginx/html:ro", BuildDirectory), "--publish", fmt.Sprintf("%d:80", port), "nginx"}...)
		if err := cmd.Run(); err != nil {
			log.Fatal(err)
		}
	}
}

func docker() string {
	bin, err := exec.LookPath("docker")
	if err != nil {
		log.Fatal("docker missing")
	}

	return bin
}
