package main

import (
	"bufio"
	"os"
	"path/filepath"

	"github.com/docker/docker/builder/dockerfile/parser"
)

func sampleFile() string {
	return settings.RepositoryDirectory() + string(filepath.Separator) + "Dockerfile"
}

func Parse() {
	dockerfile := filepath.Join(sampleFile())
	parrot.Debug("Dockerfile: ", sampleFile())

	df, err := os.Open(dockerfile)
	if err != nil {
		parrot.Error("Dockerfile missing for %s: %v", sampleFile(), err)
	}

	scanner := bufio.NewScanner(df)
	for scanner.Scan() {
		parrot.Debug(scanner.Text())
	}

	d := parser.Directive{LookingForDirectives: true}
	parser.SetEscapeToken(parser.DefaultEscapeToken, &d)

	p, err := parser.Parse(df, &d)

	if err == nil {
		parrot.Error("No error parsing broken dockerfile for", sampleFile())
	}

	parrot.Info("::", p.Value)

	df.Close()

}
