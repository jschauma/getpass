package getpass

import (
	"fmt"
	"io"
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

// TestGetPassFd opens a temporary file and verifies
// that it can read from the resulting file
// descriptor.
func TestGetpassFd(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "getpass")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte("password\n")); err != nil {
		log.Fatal(err)
	}

	if _, err := tmpfile.Seek(0, io.SeekStart); err != nil {
		log.Fatal(err)
	}

	fd := os.NewFile(uintptr(tmpfile.Fd()), "fd")
	passin := fmt.Sprintf("fd:%d", fd.Fd())
	p, err := Getpass(passin)
	if err != nil || p != "password" {
		t.Fatalf(`Getpass("%s") = %q, %v, want password, nil`, passin, p, err)
	}
}

// TestGetpassFile tests that Getpass can get a
// password from a file.
func TestGetpassFile(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "getpass")
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

// TestGetpassPass tests that Getpass with a passfrom set to 'cmd:'
// executes the given command with full shell evaluation..
func TestGetpassCmd(t *testing.T) {
	want := os.Getenv("USER")
	cmd := "echo $USER"
	p, err := Getpass("cmd:" + cmd)
	if err != nil || p != want {
		t.Fatalf(`Getpass("cmd:%s") = |%q|, %v, want |%s|, nil`, cmd, p, err, want)
	}

	want = "/dev/null"
	cmd = "ls -l /dev/null | awk '{print $NF}'"
	p, err = Getpass("cmd:" + cmd)
	if err != nil || p != want {
		t.Fatalf(`Getpass("cmd:%s") = %q, %v, want %s, nil`, cmd, p, err, want)
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
