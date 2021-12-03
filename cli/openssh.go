package cli

import (
	"bytes"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/pem"

	"golang.org/x/crypto/ssh"
)

// opensshPrivateKeyWrapper represents the private key format for OpenSSH.
type opensshPrivateKeyWrapper struct {
	CipherName  string
	KDFName     string
	KDFOptions  string
	NumKeys     uint32
	PublicKeys  []byte
	PrivateKeys []byte
}

// opensshPrivateKey represents an unencrypted private key.
type opensshPrivateKey struct {
	CheckInt1  uint32
	CheckInt2  uint32
	KeyType    string
	PublicKey  []byte
	PrivateKey []byte
	Comment    string
	Pad        []byte `ssh:"rest"`
}

const (
	// opensshAuthMagic is the protocol magic null-terminated byte string.
	opensshAuthMagic = "openssh-key-v1\x00"

	// opensshCipherBlockSize is the block size for unencrypted keys.
	opensshCipherBlockSize = 8
)

// encodePrivateKey returns the OpenSSH-encoded Ed25519 private key.
func encodePrivateKey(priv ed25519.PrivateKey) ([]byte, error) {
	var buf bytes.Buffer
	check, err := newCheckInt()
	if err != nil {
		return nil, err
	}
	privateKey := opensshPrivateKey{
		CheckInt1:  check,
		CheckInt2:  check,
		KeyType:    ssh.KeyAlgoED25519,
		PublicKey:  []byte(priv.Public().(ed25519.PublicKey)),
		PrivateKey: []byte(priv),
	}
	n := opensshCipherBlockSize - (len(ssh.Marshal(privateKey)) % opensshCipherBlockSize)
	privateKey.Pad = make([]byte, n)
	for i := 0; i < n; i++ {
		privateKey.Pad[i] = byte(i + 1)
	}
	w := opensshPrivateKeyWrapper{
		CipherName:  "none",
		KDFName:     "none",
		KDFOptions:  "",
		NumKeys:     1,
		PublicKeys:  encodePublicKey(privateKey.PublicKey),
		PrivateKeys: ssh.Marshal(privateKey),
	}
	block := &pem.Block{
		Type:  "OPENSSH PRIVATE KEY",
		Bytes: append([]byte(opensshAuthMagic), ssh.Marshal(w)...),
	}
	err = pem.Encode(&buf, block)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// encodePublicKey returns the OpenSSH-encoded public key.
func encodePublicKey(pub []byte) []byte {
	b := encodeLengthPrefix([]byte(ssh.KeyAlgoED25519))
	b = append(b, encodeLengthPrefix(pub)...)
	return b
}

// encodeLengthPrefix returns b prefixed with the uint32 length of b.
func encodeLengthPrefix(b []byte) []byte {
	n := uint32(len(b))
	length := make([]byte, 4)
	length[0] = byte(n >> 24)
	length[1] = byte(n >> 16)
	length[2] = byte(n >> 8)
	length[3] = byte(n)
	return append(length, b...)
}

// newCheckInt returns a new random uint32.
func newCheckInt() (uint32, error) {
	var rv uint32
	b := make([]byte, 4)
	_, err := rand.Read(b)
	if err != nil {
		return 0, err
	}
	rv |= uint32(b[0])
	rv |= uint32(b[1]) << 8
	rv |= uint32(b[2]) << 16
	rv |= uint32(b[3]) << 24
	return rv, nil
}
