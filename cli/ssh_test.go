package cli

import (
	"errors"
	"net"
	"path/filepath"
	"testing"

	"golang.org/x/crypto/ssh"
)

type testSSH struct{}

func newTestSSH(t *testing.T, home string, privateHostKey ssh.Signer) string {
	t.Helper()
	config := &ssh.ServerConfig{
		PublicKeyCallback: func(c ssh.ConnMetadata, pub ssh.PublicKey) (*ssh.Permissions, error) {
			filename := filepath.Join(home, username, "id_ed25519.pub")
			publicKey, err := getString(filename)
			if err != nil {
				return nil, errors.New("unauthorized")
			}
			authorizedKey, _, _, _, err := ssh.ParseAuthorizedKey([]byte(publicKey))
			if err != nil {
				return nil, errors.New("unauthorized")
			}
			if string(pub.Marshal()) == string(authorizedKey.Marshal()) {
				p := &ssh.Permissions{
					Extensions: map[string]string{
						"pubkey-fp": ssh.FingerprintSHA256(pub),
					},
				}
				return p, nil
			}
			return nil, errors.New("unauthorized")
		},
	}
	config.AddHostKey(privateHostKey)
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	addr := listener.Addr().String()
	_, port, err := net.SplitHostPort(addr)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	ss := &testSSH{}
	go ss.serve(listener, config)
	return port
}

func (s *testSSH) serve(listener net.Listener, config *ssh.ServerConfig) {
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go s.handleConnection(conn, config)
	}
}

func (s *testSSH) handleConnection(nconn net.Conn, config *ssh.ServerConfig) {
	sconn, chans, reqs, err := ssh.NewServerConn(nconn, config)
	if err != nil {
		return
	}
	defer sconn.Close()
	go ssh.DiscardRequests(reqs)
	for ch := range chans {
		go s.handleNewChannel(ch)
	}
}

func (s *testSSH) handleNewChannel(nch ssh.NewChannel) {
	if nch.ChannelType() != "session" {
		nch.Reject(ssh.UnknownChannelType, "unknown channel type")
		return
	}
	ch, reqs, err := nch.Accept()
	if err != nil {
		return
	}
	go s.handle(ch, reqs)
}

func (s *testSSH) handle(ch ssh.Channel, reqs <-chan *ssh.Request) error {
	for req := range reqs {
		switch req.Type {
		case "exec":
			var payload = struct{ Value string }{}
			err := ssh.Unmarshal(req.Payload, &payload)
			if err != nil {
				return err
			}
			switch payload.Value {
			case "docker container inspect -f {{.Id}} acroboxd":
			case "docker exec acroboxd acroboxd status":
			default:
				return req.Reply(false, nil)
			}
			err = req.Reply(true, nil)
			if err != nil {
				return err
			}
			status := struct{ Status uint32 }{uint32(0)}
			_, err = ch.SendRequest("exit-status", false, ssh.Marshal(&status))
			if err != nil {
				return err
			}
			return ch.Close()
		}
	}
	return nil
}
