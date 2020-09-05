package main

import (
	"strings"
	"testing"
)

func TestParseArgsNoPrgName(t *testing.T) {
	if err := parseArgs(nil); err != errNoProgramName {
		t.Error("Want", errNoProgramName, "got", err)
	}
}

func TestParseArgsExtraArgs(t *testing.T) {
	for _, osArgs := range [][]string{
		{"prgName", "an-arg"},
		{"prgName", "multiple", "args"},
	} {
		err := parseArgs(osArgs)
		if err == nil || !strings.HasPrefix(err.Error(), "unrecognized arguments: ") {
			t.Error("Want", "‘unrecognized arguments: ’", "got", err)
		}
	}
}

func TestParseArgsNoArgs(t *testing.T) {
	if err := parseArgs([]string{"prg-name"}); err != nil {
		t.Error(err)
	}
}
