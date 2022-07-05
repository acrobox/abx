package cli

import (
	"crypto/ed25519"
	"crypto/rand"
	"errors"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
)

func (c *client) waitForMachine(id string) (*getMachineResponse, error) {
	for n := 1; n < 30; n++ {
		time.Sleep(time.Duration(n) * time.Second)
		m, err := c.http.getMachine(id)
		if err != nil {
			return nil, err
		}
		if m.IPv4 == "" {
			continue
		}
		return m, nil
	}
	return nil, errors.New("Timeout exceeded while provisioning machine.")
}

func (c *client) waitForSSH(ipv4 string) error {
	addr := net.JoinHostPort(ipv4, c.flags.port)
	for n := 1; n < 30; n++ {
		conn, err := net.Dial("tcp", addr)
		if err != nil {
			time.Sleep(time.Duration(n) * time.Second)
			continue
		}
		return conn.Close()
	}
	return errors.New("Timeout exceeded while waiting for SSH connectivity.")
}

func (c *client) waitForAcrobox() error {
	for n := 1; n < 30; n++ {
		_, _, err := c.run("docker container inspect -f {{.Id}} acroboxd")
		if err != nil {
			time.Sleep(time.Duration(n) * time.Second)
			continue
		}
		return nil
	}
	return errors.New("Timeout exceeded while waiting for machine setup.")
}

func (c *client) waitForService() error {
	for n := 1; n < 30; n++ {
		_, _, err := c.run("docker exec acroboxd acroboxd status")
		if err != nil {
			time.Sleep(time.Duration(n) * time.Second)
			continue
		}
		return nil
	}
	return errors.New("Timeout exceeded while waiting for service setup.")
}

func (c *client) writeBytes(name string, b []byte, perm os.FileMode) error {
	filename := filepath.Join(c.config.Home, c.flags.host, name)
	return os.WriteFile(filename, b, perm)
}

func (c *client) writeKey(name string, b []byte) error {
	return c.writeBytes(name, b, 0600)
}

func (c *client) writeString(name, value string) error {
	return c.writeBytes(name, []byte(value+"\n"), 0660)
}

func (c *client) getID() (string, error) {
	filename := filepath.Join(c.config.Home, c.flags.host, "ID")
	return getString(filename)
}

func (c *client) getIPv4() (string, error) {
	filename := filepath.Join(c.config.Home, c.flags.host, "IPv4")
	return getString(filename)
}

func (c *client) getKnownHosts() (ssh.HostKeyCallback, error) {
	filename := filepath.Join(c.config.Home, c.flags.host, "known_hosts")
	return knownhosts.New(filename)
}

func (c *client) getPrivateKey() (ssh.Signer, error) {
	filename := filepath.Join(c.config.Home, c.flags.host, "id_ed25519")
	b, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return ssh.ParsePrivateKey(b)
}

// newKeyPair returns a OpenSSH-encoded Ed25519 private key
// and its corresponding public SSH authorized key.
func newKeyPair() ([]byte, []byte, error) {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, nil, err
	}
	privateKey, err := encodePrivateKey(priv)
	if err != nil {
		return nil, nil, err
	}
	publicKey, err := ssh.NewPublicKey(pub)
	if err != nil {
		return nil, nil, err
	}
	authorizedKey := ssh.MarshalAuthorizedKey(publicKey)
	return privateKey, authorizedKey, nil
}

// getString returns a space-trimmed string from the filename.
func getString(filename string) (string, error) {
	b, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(b)), nil
}
