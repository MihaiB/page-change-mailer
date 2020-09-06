package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"time"
)

const (
	env_url             = "URL"
	env_filename        = "FILE"
	env_delay_min       = "DELAY_MIN"
	env_delay_max       = "DELAY_MAX"
	env_smtps_host      = "SMTPS_HOST"
	env_smtps_port      = "SMTPS_PORT"
	env_smtps_username  = "SMTPS_USERNAME"
	env_smtps_password  = "SMTPS_PASSWORD"
	env_email_addr_from = "EMAIL_ADDRESS_FROM"
	env_email_addr_to   = "EMAIL_ADDRESS_TO"
)

var env_var_names = map[string]struct{}{
	env_url:             {},
	env_filename:        {},
	env_delay_min:       {},
	env_delay_max:       {},
	env_smtps_host:      {},
	env_smtps_port:      {},
	env_smtps_username:  {},
	env_smtps_password:  {},
	env_email_addr_from: {},
	env_email_addr_to:   {},
}

func getEnv() map[string]string {
	env := map[string]string{}
	for k := range env_var_names {
		env[k] = os.Getenv(k)
	}
	return env
}

var errNoProgramName = errors.New("no program name (os.Args empty)")

type argsT struct {
	url, filename                  string
	delay_min, delay_max           time.Duration
	smtps_host, smtps_port         string
	smtps_username, smtps_password string
	email_addr_from, email_addr_to string
}

func parseEnv(env map[string]string) (*argsT, error) {
	args := &argsT{}

	if args.url = env[env_url]; args.url == "" {
		return nil, errors.New("empty env var: " + env_url)
	}
	if args.filename = env[env_filename]; args.filename == "" {
		return nil, errors.New("empty env var: " + env_filename)
	}

	var err error
	if args.delay_min, err = time.ParseDuration(env[env_delay_min]); err != nil {
		return nil, err
	}
	if args.delay_max, err = time.ParseDuration(env[env_delay_max]); err != nil {
		return nil, err
	}

	if args.delay_min < 0 {
		return nil, fmt.Errorf("negative %v: %v",
			env_delay_min, args.delay_min)
	}
	if args.delay_min > args.delay_max {
		return nil, fmt.Errorf("%v %v > %v %v",
			env_delay_min, args.delay_min,
			env_delay_max, args.delay_max)
	}

	if args.smtps_host = env[env_smtps_host]; args.smtps_host == "" {
		return nil, errors.New("empty env var: " + env_smtps_host)
	}
	if args.smtps_port = env[env_smtps_port]; args.smtps_port == "" {
		return nil, errors.New("empty env var: " + env_smtps_port)
	}
	if args.smtps_username = env[env_smtps_username]; args.smtps_username == "" {
		return nil, errors.New("empty env var: " + env_smtps_username)
	}
	if args.smtps_password = env[env_smtps_password]; args.smtps_password == "" {
		return nil, errors.New("empty env var: " + env_smtps_password)
	}
	if args.email_addr_from = env[env_email_addr_from]; args.email_addr_from == "" {
		return nil, errors.New("empty env var: " + env_email_addr_from)
	}
	if args.email_addr_to = env[env_email_addr_to]; args.email_addr_to == "" {
		return nil, errors.New("empty env var: " + env_email_addr_to)
	}

	return args, nil
}

func parseArgs(osArgs []string, env map[string]string) (*argsT, error) {
	if len(osArgs) == 0 {
		return nil, errNoProgramName
	}

	fs := flag.NewFlagSet(osArgs[0], flag.ExitOnError)
	fs.Usage = func() {
		fmt.Fprint(fs.Output(), "usage: ", fs.Name(), "\n")
		fmt.Fprint(fs.Output(), `
Fetch a webpage regularly and send an email when it changes.

See the README for how to set parameters.
`)
		fs.PrintDefaults()
	}

	if err := fs.Parse(osArgs[1:]); err != nil {
		return nil, err
	}
	if fs.NArg() > 0 {
		return nil, fmt.Errorf("unrecognized arguments: %#v", fs.Args())
	}

	return parseEnv(env)
}
