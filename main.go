package main

import (
	"bytes"
	"crypto/tls"
	"github.com/jordan-wright/email"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/smtp"
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

func sendEmail(args *argsT, newContent []byte) error {
	e := email.NewEmail()
	e.From = args.email_addr_from
	e.To = []string{args.email_addr_to}
	e.Subject = "Page changed: " + args.url
	e.Text = []byte(args.url)
	e.HTML = newContent

	addr := args.smtps_host + ":" + args.smtps_port
	auth := smtp.PlainAuth("", args.smtps_username, args.smtps_password,
		args.smtps_host)
	tlsConfig := &tls.Config{ServerName: args.smtps_host}

	if err := e.SendWithTLS(addr, auth, tlsConfig); err != nil {
		return err
	}

	logger.Println("email sent from", args.email_addr_from,
		"to", args.email_addr_to)
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
		err = sendEmail(args, newContent)
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
