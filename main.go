package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

func init() {
	rand.Seed(time.Now().Unix())
}

var logger = log.New(os.Stderr, "UTC ", log.LstdFlags|log.LUTC|log.Lmsgprefix)

func fetchPage(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func shouldEmail(filename string, newContent []byte) (bool, error) {
	oldContent, err := ioutil.ReadFile(filename)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return !bytes.Equal(oldContent, newContent), nil
}

func sendEmail() error {
	logger.Print("TODO: send email")
	return nil
}

func fetchAndEmail(args *argsT) error {
	newContent, err := fetchPage(args.url)
	if err != nil {
		return err
	}

	notify, err := shouldEmail(args.filename, newContent)
	if err != nil {
		return err
	}

	if notify {
		err = sendEmail()
		if err != nil {
			return err
		}
	}

	return ioutil.WriteFile(args.filename, newContent, 0644)
}

func sleep(delay_min, delay_max time.Duration) {
	d := delay_min + time.Duration(rand.Int63n(int64(delay_max-delay_min+1)))
	logger.Println("sleeping for", d.Round(time.Second))
	time.Sleep(d)
}

func main_err() error {
	args, err := parseArgs(os.Args, getEnv())
	if err != nil {
		return err
	}

	for {
		if err = fetchAndEmail(args); err != nil {
			logger.Print(err)
		}
		sleep(args.delay_min, args.delay_max)
	}
}

func main() {
	if err := main_err(); err != nil {
		logger.Fatal(err)
	}
}
