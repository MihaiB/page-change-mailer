package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"time"
)

const (
	env_url       = "URL"
	env_filename  = "FILE"
	env_delay_min = "DELAY_MIN"
	env_delay_max = "DELAY_MAX"
)

var env_var_names = []string{env_url, env_filename, env_delay_min, env_delay_max}

func getEnv() map[string]string {
	env := map[string]string{}
	for _, k := range env_var_names {
		env[k] = os.Getenv(k)
	}
	return env
}

var errNoProgramName = errors.New("no program name (os.Args empty)")

type argsT struct {
	url, filename        string
	delay_min, delay_max time.Duration
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
