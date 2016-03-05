package reloader

import (
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
	"time"
)

func noop([]byte, error) error {
	return nil
}

func TestNoStat(t *testing.T) {
	filename := os.TempDir() + "/doesntexist.123456789"
	err := New(filename, noop)
	if err == nil {
		t.Errorf("Expected New to return error when the file doesn't exist.")
	}
}

func TestNoRead(t *testing.T) {
	// Create a file with no permissions.
	filename := os.TempDir() + "/test-no-read.txt"
	ioutil.WriteFile(filename, []byte{}, 0)
	err := New(filename, noop)
	if err == nil {
		t.Errorf("Expected New to return error when permission denied.")
	}
}

func TestFirstError(t *testing.T) {
	filename := os.TempDir() + "/test-first-error.txt"
	ioutil.WriteFile(filename, []byte{}, 0644)
	err := New(filename, func([]byte, error) error {
		return fmt.Errorf("i die")
	})
	if err == nil {
		t.Errorf("Expected New to return error when the callback returned error the first time.")
	}
}

func TestFirstSuccess(t *testing.T) {
	filename := os.TempDir() + "/test-first-success.txt"
	ioutil.WriteFile(filename, []byte{}, 0644)
	err := New(filename, func([]byte, error) error {
		return nil
	})
	if err != nil {
		t.Errorf("Expected New to succeed.")
	}
}

func TestReload(t *testing.T) {
	filename := os.TempDir() + "/test-reload.txt"
	ioutil.WriteFile(filename, []byte("first body"), 0644)

	var bodies []string
	reloads := make(chan []byte, 1)
	err := New(filename, func(b []byte, err error) error {
		if err != nil {
			t.Fatalf("Got error in callback: %s", err)
		}
		bodies = append(bodies, string(b))
		reloads <- b
		return nil
	})
	if err != nil {
		t.Fatalf("Expected New to succeed.")
	}
	<-reloads
	expected := []string{"first body"}
	if !reflect.DeepEqual(bodies, expected) {
		t.Errorf("Expected bodies = %#v, got %#v", expected, bodies)
	}
	time.Sleep(2 * time.Second)
	if !reflect.DeepEqual(bodies, expected) {
		t.Errorf("Expected bodies = %#v, got %#v", expected, bodies)
	}

	// Write to the file, expect a reload.
	ioutil.WriteFile(filename, []byte("second body"), 0644)
	time.Sleep(2 * time.Second)
	select {
	case <-reloads:
	case <-time.After(5 * time.Second):
		t.Fatalf("timed out waiting for reload")
	}
	expected = []string{"first body", "second body"}
	if !reflect.DeepEqual(bodies, expected) {
		t.Errorf("Expected bodies = %#v, got %#v", expected, bodies)
	}

	time.Sleep(2 * time.Second)
	if !reflect.DeepEqual(bodies, expected) {
		t.Errorf("Expected bodies = %#v, got %#v", expected, bodies)
	}
}

func TestReloadFailure(t *testing.T) {
	filename := os.TempDir() + "/test-reload.txt"
	ioutil.WriteFile(filename, []byte("first body"), 0644)

	type res struct {
		b   []byte
		err error
	}

	reloads := make(chan res, 1)
	err := New(filename, func(b []byte, err error) error {
		reloads <- res{b, err}
		return nil
	})
	if err != nil {
		t.Fatalf("Expected New to succeed.")
	}
	<-reloads
	os.Remove(filename)
	time.Sleep(2 * time.Second)
	select {
	case r := <-reloads:
		if r.err == nil {
			t.Errorf("Expected error trying to read missing file.")
		}
	case <-time.After(5 * time.Second):
		t.Fatalf("timed out waiting for reload")
	}

	// Create a file with no permissions
	ioutil.WriteFile(filename, []byte("second body"), 0)
	time.Sleep(2 * time.Second)
	select {
	case r := <-reloads:
		if r.err == nil {
			t.Errorf("Expected error trying to read file with no permissions.")
		}
	case <-time.After(5 * time.Second):
		t.Fatalf("timed out waiting for reload")
	}

	err = os.Remove(filename)
	if err != nil {
		t.Fatal(err)
	}
	err = ioutil.WriteFile(filename, []byte("third body"), 0644)
	if err != nil {
		t.Fatal(err)
	}
	for {
		select {
		case r := <-reloads:
			if r.err != nil {
				continue
			}
			if string(r.b) != "third body" {
				t.Errorf("Expected 'third body' reading file after restoring it.")
			}
			return
		case <-time.After(5 * time.Second):
			t.Fatalf("timed out waiting for successful reload")
		}
	}
}
