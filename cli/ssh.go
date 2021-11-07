package cli

import (
	"bytes"
	"io"
	"net"
	"os"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"

	"github.com/pnelson/cli"
)

const term = "xterm-256color"

func (c *client) newSession() (*ssh.Session, error) {
	ipv4, err := c.getIPv4()
	if err != nil {
		return nil, err
	}
	privateKey, err := c.getPrivateKey()
	if err != nil {
		return nil, err
	}
	knownHosts, err := c.getKnownHosts()
	if err != nil {
		return nil, err
	}
	config := &ssh.ClientConfig{
		User:            username,
		Auth:            []ssh.AuthMethod{ssh.PublicKeys(privateKey)},
		HostKeyCallback: knownHosts,
		Timeout:         90 * time.Second,
	}
	addr := net.JoinHostPort(ipv4, c.flags.port)
	s, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		return nil, err
	}
	return s.NewSession()
}

func (c *client) run(command string) ([]byte, []byte, error) {
	return c.runWithStdin(nil, command)
}

func (c *client) runWithStdin(stdin io.Reader, command string) ([]byte, []byte, error) {
	session, err := c.newSession()
	if err != nil {
		return nil, nil, err
	}
	defer session.Close()
	stdout := bytes.Buffer{}
	stderr := bytes.Buffer{}
	session.Stdin = stdin
	session.Stdout = &stdout
	session.Stderr = &stderr
	err = session.Run(command)
	if err != nil {
		return nil, nil, err
	}
	return stdout.Bytes(), stderr.Bytes(), nil
}

func (c *client) exec(command string, args []string) error {
	session, err := c.newSession()
	if err != nil {
		return err
	}
	defer session.Close()
	session.Stdin = c.config.Stdin
	session.Stdout = c.config.Stdout
	session.Stderr = c.config.Stderr
	modes := ssh.TerminalModes{
		ssh.ECHO:          0,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}
	f, ok := session.Stdin.(*os.File)
	if ok {
		fd := int(f.Fd())
		if terminal.IsTerminal(fd) {
			state, err := terminal.MakeRaw(fd)
			if err != nil {
				return err
			}
			defer terminal.Restore(fd, state)
			width, height, err := terminal.GetSize(fd)
			if err != nil {
				return err
			}
			err = session.RequestPty(term, height, width, modes)
			if err != nil {
				return err
			}
		}
	}
	// Arguments will have already been processed by the shell
	// before they are interpreted by abx. We need to re-quote
	// the arguments for the remote shell to interpret them as
	// they were initially passed. For example:
	//
	//   $ abx env/set image key="value with spaces"
	//
	// This args slice will contain the string "key=value with spaces"
	// but if proxying the command directly to the remote server then
	// it would only see "key=value" and then bail out attempting to
	// process the next key/value pair.
	for i, arg := range args {
		args[i] = "'" + strings.ReplaceAll(arg, "'", "'\"'\"'") + "'"
	}
	command = command + " " + strings.Join(args, " ")
	command = strings.TrimSpace(command) + "\n"
	return session.Run(command)
}

func (c *client) acroboxd(command string, args []string) error {
	err := c.exec("docker exec -i -t acroboxd acroboxd "+command, args)
	_, ok := err.(*ssh.ExitError)
	if ok {
		return cli.ErrExitFailure
	}
	return err
}
