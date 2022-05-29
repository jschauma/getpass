package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/jschauma/getpass"
)

func main() {

	var pass string
	flag.StringVar(&pass, "p", "tty", "password method")
	flag.Parse()

	// Try out any of the options:
	p, err := getpass.Getpass(pass)
	if err != nil {
		log.Fatal("Unable to get password from user: ", err)
	} else {
		fmt.Printf("%s\n", p)
	}

	// Alternatively:

	// This will prompt the user to enter a password interactively,
	// using the default prompt.
	p, err = getpass.Getpass()
	if err != nil {
		log.Fatal("Unable to get password from user: %v\n", err)
	} else {
		fmt.Printf("%s\n", p)
	}

	// Using a custom prompt:
	p, err = getpass.Getpass("tty:Please enter your secret passphrase: ")
	if err != nil {
		log.Fatal("Unable to get password from user: %v\n", err)
	} else {
		fmt.Printf("%s\n", p)
	}

	// Using an environment variable:
	p, err = getpass.Getpass("env:MYSECRET")
	if err != nil {
		log.Fatal("Unable to get password from user: ", err)
	} else {
		fmt.Printf("%s\n", p)
	}

	// Using a file:
	p, err = getpass.Getpass("file:~/.secret")
	if err != nil {
		log.Fatal("Unable to get password from user: ", err)
	} else {
		fmt.Printf("%s\n", p)
	}

	// Using a command:
	p, err = getpass.Getpass("cmd:env | grep -i mysecret")
	if err != nil {
		log.Fatal("Unable to get password from user: ", err)
	} else {
		fmt.Printf("%s\n", p)
	}
	// etc. etc.
}
