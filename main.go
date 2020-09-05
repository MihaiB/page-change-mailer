package main

import (
	"fmt"
	"log"
	"os"
)

var logger = log.New(os.Stderr, "UTC ", log.LstdFlags|log.LUTC|log.Lmsgprefix)

func main_err() error {
	args, err := parseArgs(os.Args, getEnv())
	if err != nil {
		return err
	}

	fmt.Println(args)

	return nil
}

func main() {
	if err := main_err(); err != nil {
		logger.Fatal(err)
	}
}
