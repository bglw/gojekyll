package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	libsass "github.com/wellington/go-libsass"
)

func (p *DynamicPage) writeSass(w io.Writer, data []byte) error {
	comp, err := libsass.New(w, bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	err = comp.Option(libsass.IncludePaths(site.SassIncludePaths()))
	if err != nil {
		log.Fatal(err)
	}
	return comp.Run()
}

// CopySassFileIncludes copies sass partials into a temporary directory,
// removing initial underscores.
// TODO delete the temp directory when done
func (s *Site) CopySassFileIncludes() {
	// TODO use libsass.ImportsOption instead?
	if site.sassTempDir == "" {
		d, err := ioutil.TempDir(os.TempDir(), "_sass")
		if err != nil {
			panic(err)
		}
		site.sassTempDir = d
	}

	src := filepath.Join(s.Source, "_sass")
	dst := site.sassTempDir
	err := filepath.Walk(src, func(from string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}
		rel, err := filepath.Rel(src, from)
		if err != nil {
			panic(err)
		}
		if strings.HasPrefix(rel, "_") {
			rel = rel[1:]
		}
		to := filepath.Join(dst, rel)
		if err := copyFile(to, from, 0644); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
}

// SassIncludePaths returns an array of sass include directories.
func (s *Site) SassIncludePaths() []string {
	if site.sassTempDir == "" {
		site.CopySassFileIncludes()
	}
	s.CopySassFileIncludes()
	return []string{site.sassTempDir}
}