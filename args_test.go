package main

import (
	"strings"
	"testing"
)

// Get a new map every time before passing it to the code under test
// because the map is mutable.
func getValidEnv() map[string]string {
	return map[string]string{
		envURL:           "https://example.com",
		envFilename:      "some/path/to/file",
		envDelayMin:      "14m30s",
		envDelayMax:      "24h",
		envSMTPSHost:     "example.com",
		envSMTPSPort:     "465",
		envSMTPSUsername: "user",
		envSMTPSPassword: "pass",
		envEmailAddrFrom: "from@example.com",
		envEmailAddrTo:   "to@example.com",
	}
}

func TestParseArgsNoPrgName(t *testing.T) {
	if _, err := parseArgs(nil, getValidEnv()); err != errNoProgramName {
		t.Error("Want", errNoProgramName, "got", err)
	}
}

func TestParseArgsExtraArgs(t *testing.T) {
	for _, osArgs := range [][]string{
		{"prgName", "an-arg"},
		{"prgName", "multiple", "args"},
	} {
		_, err := parseArgs(osArgs, getValidEnv())
		if err == nil || !strings.HasPrefix(err.Error(), "unrecognized arguments: ") {
			t.Error("Want", "‘unrecognized arguments: ’", "got", err)
		}
	}
}

func TestParseArgsNoArgs(t *testing.T) {
	args, err := parseArgs([]string{"prg-name"}, getValidEnv())
	if err != nil {
		t.Error(err)
	}
	if args == nil {
		t.Error("nil args")
	}
}

func TestParseEnvInvalid(t *testing.T) {
	for _, tc := range []*struct {
		k, v string
		want string
	}{
		{envURL, "", "empty env var: " + envURL},
		{envFilename, "", "empty env var: " + envFilename},
		{envDelayMin, "", `time: invalid duration ""`},
		{envDelayMax, "17", `time: missing unit in duration "17"`},

		{envDelayMax, "5m", "DELAY_MIN 14m30s > DELAY_MAX 5m0s"},
		{envDelayMin, "-37s", "negative DELAY_MIN: -37s"},

		{envSMTPSHost, "", "empty env var: " + envSMTPSHost},
		{envSMTPSPort, "", "empty env var: " + envSMTPSPort},
		{envSMTPSUsername, "", "empty env var: " + envSMTPSUsername},
		{envSMTPSPassword, "", "empty env var: " + envSMTPSPassword},
		{envEmailAddrFrom, "", "empty env var: " + envEmailAddrFrom},
		{envEmailAddrTo, "", "empty env var: " + envEmailAddrTo},
	} {
		e := getValidEnv()
		e[tc.k] = tc.v
		args, err := parseEnv(e)
		if err == nil || err.Error() != tc.want {
			t.Error("Want", tc.want, "got", err)
		}
		if args != nil {
			t.Error("Got non-nil args:", args)
		}
	}

}
