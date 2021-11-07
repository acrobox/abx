package cli

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io/ioutil"
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
	return ioutil.WriteFile(filename, b, perm)
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
	filename := filepath.Join(c.config.Home, c.flags.host, "id_rsa")
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return ssh.ParsePrivateKey(b)
}

// newKeyPair returns a PEM-encoded RSA private key
// and its corresponding public SSH authorized key.
func newKeyPair() ([]byte, []byte, error) {
	priv, err := rsa.GenerateKey(rand.Reader, keySize)
	if err != nil {
		return nil, nil, err
	}
	privateKey, err := encodePrivateKey(priv)
	if err != nil {
		return nil, nil, err
	}
	pub, err := ssh.NewPublicKey(&priv.PublicKey)
	if err != nil {
		return nil, nil, err
	}
	authorizedKey := ssh.MarshalAuthorizedKey(pub)
	return privateKey, authorizedKey, nil
}

// encodePrivateKey returns the PEM-encoded RSA private key.
func encodePrivateKey(priv *rsa.PrivateKey) ([]byte, error) {
	var buf bytes.Buffer
	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(priv),
	}
	err := pem.Encode(&buf, block)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// getString returns a space-trimmed string from the filename.
func getString(filename string) (string, error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(b)), nil
}
