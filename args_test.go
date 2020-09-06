package main

import (
	"strings"
	"testing"
)

// Get a new map every time before passing it to the code under test
// because the map is mutable.
func getValidEnv() map[string]string {
	return map[string]string{
		env_url:             "https://example.com",
		env_filename:        "some/path/to/file",
		env_delay_min:       "14m30s",
		env_delay_max:       "24h",
		env_smtps_host:      "example.com",
		env_smtps_port:      "465",
		env_smtps_username:  "user",
		env_smtps_password:  "pass",
		env_email_addr_from: "from@example.com",
		env_email_addr_to:   "to@example.com",
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
		{env_url, "", "empty env var: " + env_url},
		{env_filename, "", "empty env var: " + env_filename},
		{env_delay_min, "", `time: invalid duration ""`},
		{env_delay_max, "17", `time: missing unit in duration "17"`},

		{env_delay_max, "5m", "DELAY_MIN 14m30s > DELAY_MAX 5m0s"},
		{env_delay_min, "-37s", "negative DELAY_MIN: -37s"},

		{env_smtps_host, "", "empty env var: " + env_smtps_host},
		{env_smtps_port, "", "empty env var: " + env_smtps_port},
		{env_smtps_username, "", "empty env var: " + env_smtps_username},
		{env_smtps_password, "", "empty env var: " + env_smtps_password},
		{env_email_addr_from, "", "empty env var: " + env_email_addr_from},
		{env_email_addr_to, "", "empty env var: " + env_email_addr_to},
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
