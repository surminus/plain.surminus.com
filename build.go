package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func build() error {
	if err := os.RemoveAll(BuildDirectory); err != nil {
		return err
	}
	if err := os.Mkdir(BuildDirectory, 0755); err != nil {
		return err
	}

	var stylesheet string

	// If styles.css exists, copy it into the build directory and ensure it's
	// included in the generation
	if f, err := os.Stat("styles.css"); err == nil {
		data, err := os.ReadFile(f.Name())
		if err != nil {
			return err
		}

		if err := os.WriteFile(filepath.Join(BuildDirectory, "styles.css"), data, 0644); err != nil {
			return err
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
		return err
	}

	return nil
}
