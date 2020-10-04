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

var envVarNames = map[string]struct{}{
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
	for k := range envVarNames {
		env[k] = os.Getenv(k)
	}
	return env
}

var errNoProgramName = errors.New("no program name (os.Args empty)")

type argsT struct {
	url, filename                string
	delayMin, delayMax           time.Duration
	smtpsHost, smtpsPort         string
	smtpsUsername, smtpsPassword string
	emailAddrFrom, emailAddrTo   string
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
	if args.delayMin, err = time.ParseDuration(env[envDelayMin]); err != nil {
		return nil, err
	}
	if args.delayMax, err = time.ParseDuration(env[envDelayMax]); err != nil {
		return nil, err
	}

	if args.delayMin < 0 {
		return nil, fmt.Errorf("negative %v: %v",
			envDelayMin, args.delayMin)
	}
	if args.delayMin > args.delayMax {
		return nil, fmt.Errorf("%v %v > %v %v",
			envDelayMin, args.delayMin,
			envDelayMax, args.delayMax)
	}

	if args.smtpsHost = env[envSMTPSHost]; args.smtpsHost == "" {
		return nil, errors.New("empty env var: " + envSMTPSHost)
	}
	if args.smtpsPort = env[envSMTPSPort]; args.smtpsPort == "" {
		return nil, errors.New("empty env var: " + envSMTPSPort)
	}
	if args.smtpsUsername = env[envSMTPSUsername]; args.smtpsUsername == "" {
		return nil, errors.New("empty env var: " + envSMTPSUsername)
	}
	if args.smtpsPassword = env[envSMTPSPassword]; args.smtpsPassword == "" {
		return nil, errors.New("empty env var: " + envSMTPSPassword)
	}
	if args.emailAddrFrom = env[envEmailAddrFrom]; args.emailAddrFrom == "" {
		return nil, errors.New("empty env var: " + envEmailAddrFrom)
	}
	if args.emailAddrTo = env[envEmailAddrTo]; args.emailAddrTo == "" {
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
