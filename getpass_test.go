package getpass

import (
	"io/ioutil"
	"log"
	"os"
	"testing"
)

// TestGetpassEnv sets an environment variable and
// then tests that Getpass with a passfrom set to
// 'env:GETPASS' correctly returns it.
func TestGetpassEnv(t *testing.T) {
	want := "password"
	_ = os.Setenv("GETPASS", want)
	p, err := Getpass("env:GETPASS")
	if err != nil || p != want {
		t.Fatalf(`Getpass("env:GETPASS") = %q, %v, want %s, nil`, p, err, want)
	}
}

// TestGetpassFile tests that Getpass can get a
// password from a file.
func TestGetpassFile(t *testing.T) {
	tmpfile, err := ioutil.TempFile("", "getpass")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte("password\n")); err != nil {
		log.Fatal(err)
	}

	if err := tmpfile.Close(); err != nil {
		log.Fatal(err)
	}

	passin := "file:" + tmpfile.Name()
	p, err := Getpass(passin)
	if err != nil || p != "password" {
		t.Fatalf(`Getpass("%s") = %q, %v, want password, nil`, passin, p, err)
	}
}

// TestGetpassPass tests that Getpass with a passfrom
// set to 'pass:password' returns 'password'.
func TestGetpassPass(t *testing.T) {
	p, err := Getpass("pass:password")
	if err != nil || p != "password" {
		t.Fatalf(`Getpass("pass:password") = %q, %v, want password, nil`, p, err)
	}
}

// TestGetpassFail tests that Getpass with an invalid
// passfrom fails.
func TestGetpassFail(t *testing.T) {
	p, err := Getpass("whatever:invalid")
	if err == nil || p != "" {
		t.Fatalf(`Getpass("whatever:invalid") = %q, %v, want "", error`, p, err)
	}
}
