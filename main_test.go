package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestShouldEmailNoHistory(t *testing.T) {
	d := t.TempDir()
	filename := filepath.Join(d, "no-such-file")
	for _, s := range []string{"", "\n", "hello"} {
		val, err := shouldEmail(filename, []byte(s))
		if err != nil {
			t.Error(err)
		}
		if val {
			t.Error("wants to email when there is no " +
				"previous version of the page")
		}
	}
}

func TestShouldEmailError(t *testing.T) {
	d := t.TempDir()
	if _, err := shouldEmail(d, nil); err == nil {
		t.Error("Expected an error when the filename is a directory")
	}
}

func TestShouldEmailCompareContent(t *testing.T) {
	d := t.TempDir()
	f := filepath.Join(d, "file")
	oldContent, newContent := []byte("old"), []byte("new")
	if err := os.WriteFile(f, oldContent, 0644); err != nil {
		t.Fatal(err)
	}

	if val, err := shouldEmail(f, oldContent); val || err != nil {
		t.Error(val, err)
	}
	if val, err := shouldEmail(f, newContent); !val || err != nil {
		t.Error(val, err)
	}
}
