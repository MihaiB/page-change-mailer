package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"time"
)

const (
	envURL           = "URL"
	envFilename      = "FILE"
	envDelayMin      = "DELAY_MIN"
	envDelayMax      = "DELAY_MAX"
	envSMTPSHost     = "SMTPS_HOST"
	envSMTPSPort     = "SMTPS_PORT"
	envSMTPSUsername = "SMTPS_USERNAME"
	envSMTPSPassword = "SMTPS_PASSWORD"
	envEmailAddrFrom = "EMAIL_ADDRESS_FROM"
	envEmailAddrTo   = "EMAIL_ADDRESS_TO"
)

var env_var_names = map[string]struct{}{
	envURL:           {},
	envFilename:      {},
	envDelayMin:      {},
	envDelayMax:      {},
	envSMTPSHost:     {},
	envSMTPSPort:     {},
	envSMTPSUsername: {},
	envSMTPSPassword: {},
	envEmailAddrFrom: {},
	envEmailAddrTo:   {},
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

	if args.url = env[envURL]; args.url == "" {
		return nil, errors.New("empty env var: " + envURL)
	}
	if args.filename = env[envFilename]; args.filename == "" {
		return nil, errors.New("empty env var: " + envFilename)
	}

	var err error
	if args.delay_min, err = time.ParseDuration(env[envDelayMin]); err != nil {
		return nil, err
	}
	if args.delay_max, err = time.ParseDuration(env[envDelayMax]); err != nil {
		return nil, err
	}

	if args.delay_min < 0 {
		return nil, fmt.Errorf("negative %v: %v",
			envDelayMin, args.delay_min)
	}
	if args.delay_min > args.delay_max {
		return nil, fmt.Errorf("%v %v > %v %v",
			envDelayMin, args.delay_min,
			envDelayMax, args.delay_max)
	}

	if args.smtps_host = env[envSMTPSHost]; args.smtps_host == "" {
		return nil, errors.New("empty env var: " + envSMTPSHost)
	}
	if args.smtps_port = env[envSMTPSPort]; args.smtps_port == "" {
		return nil, errors.New("empty env var: " + envSMTPSPort)
	}
	if args.smtps_username = env[envSMTPSUsername]; args.smtps_username == "" {
		return nil, errors.New("empty env var: " + envSMTPSUsername)
	}
	if args.smtps_password = env[envSMTPSPassword]; args.smtps_password == "" {
		return nil, errors.New("empty env var: " + envSMTPSPassword)
	}
	if args.email_addr_from = env[envEmailAddrFrom]; args.email_addr_from == "" {
		return nil, errors.New("empty env var: " + envEmailAddrFrom)
	}
	if args.email_addr_to = env[envEmailAddrTo]; args.email_addr_to == "" {
		return nil, errors.New("empty env var: " + envEmailAddrTo)
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
