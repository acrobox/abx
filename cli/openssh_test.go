package cli

import (
	"bytes"
	"crypto/ed25519"
	"crypto/rand"
	"reflect"
	"testing"

	"golang.org/x/crypto/ssh"
)

func TestEncodePrivateKey(t *testing.T) {
	_, want, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	priv, err := encodePrivateKey(want)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	key, err := ssh.ParseRawPrivateKey(priv)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	have, ok := key.(*ed25519.PrivateKey)
	if !ok {
		t.Fatalf("ParseRawPrivateKey should be a pointer to ed25519.PrivateKey\nhave %#v", have)
	}
	if !reflect.DeepEqual(*have, want) {
		t.Fatalf("ParseRawPrivateKey\nhave '%s'\nwant '%s'", *have, want)
	}
}

func TestEncodeLengthPrefix(t *testing.T) {
	b := []byte(ssh.KeyAlgoED25519)
	want := append([]byte{0x00, 0x00, 0x00, 0x0b}, b...)
	have := encodeLengthPrefix(b)
	if !bytes.Equal(have, want) {
		t.Fatalf("encodeLengthPrefix\nhave '%x'\nwant '%x'", have, want)
	}
}
