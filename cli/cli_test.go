package cli

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"golang.org/x/crypto/ssh"

	"github.com/pnelson/cli"
)

func TestInit(t *testing.T) {
	privateHostKey, publicHostKey := newTestHostKeyPair(t)
	view := getMachineResponse{
		ID:        "test",
		UserID:    "test-user",
		Name:      "acrobox",
		IPv4:      "127.0.0.1",
		PublicKey: publicHostKey,
		CreatedAt: time.Now().UTC(),
	}
	fn := newTestHandler(t, http.StatusCreated, view)
	ts := httptest.NewServer(fn)
	defer ts.Close()
	home := t.TempDir()
	port := newTestSSH(t, home, privateHostKey)
	config := &Config{
		Args:   []string{"abx", "-addr", ts.URL, "-port", port, "init", "-force"},
		Home:   home,
		Stdout: ioutil.Discard,
		Stderr: ioutil.Discard,
	}
	err := Run(config)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	dir := filepath.Join(config.Home, username)
	_, err = os.Stat(dir)
	if os.IsNotExist(err) {
		t.Fatalf("host data directory '%s' should exist", dir)
	}
	filenames := []string{"IPv4", "known_hosts", "id_rsa", "id_rsa.pub"}
	for _, name := range filenames {
		filename := filepath.Join(dir, name)
		_, err = os.Stat(filename)
		if os.IsNotExist(err) {
			t.Errorf("data file '%s' should exist", filename)
		}
	}
}

func TestInitErr(t *testing.T) {
	view := errorResponse{
		Code:      http.StatusTeapot,
		Title:     http.StatusText(http.StatusTeapot),
		Message:   "test",
		RequestID: "test",
	}
	fn := newTestHandler(t, http.StatusTeapot, view)
	ts := httptest.NewServer(fn)
	defer ts.Close()
	config := &Config{
		Args:   []string{"abx", "-addr", ts.URL, "init", "-force"},
		Home:   t.TempDir(),
		Stdout: ioutil.Discard,
		Stderr: ioutil.Discard,
	}
	err := Run(config)
	if err != cli.ErrExitFailure {
		t.Fatalf("unexpected error: %v", err)
	}
}

func newTestHostKeyPair(t *testing.T) (ssh.Signer, string) {
	t.Helper()
	privateKey, authorizedKey, err := newKeyPair()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	priv, err := ssh.ParsePrivateKey(privateKey)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	return priv, string(authorizedKey)
}

func newTestHandler(t *testing.T, code int, view interface{}) http.Handler {
	t.Helper()
	fn := func(w http.ResponseWriter, req *http.Request) {
		b, err := json.Marshal(view)
		if err != nil {
			t.Fatalf("unexpected err: %v", err)
		}
		w.WriteHeader(code)
		w.Write(b)
	}
	return http.HandlerFunc(fn)
}
