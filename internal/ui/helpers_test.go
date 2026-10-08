package ui

import (
	"io"
	"os"
	"strings"
)

func newStringReader(s string) io.Reader { return strings.NewReader(s) }

func writeFileHelper(path, body string) error { return os.WriteFile(path, []byte(body), 0o644) }

func openFileHelper(path string) (*os.File, error) { return os.Open(path) }
