package cli // import "acrobox.io/abx/cli"

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
	"golang.org/x/crypto/ssh/terminal"

	"github.com/pnelson/cli"

	"acrobox.io/docs"
)

// client represents the central manager of application activity.
type client struct {
	config *Config
	cli    *cli.CLI
	http   *service
	flags  flags
}

// Config represents the core configuration parameters.
type Config struct {
	Args   []string
	Home   string
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

// Application constants.
const (
	AppName    = "abx"
	AppVersion = "0.0.1"
)

// username is the machine username and default machine name.
const username = "acrobox"

// cardAuthText is the authorization text to be displayed for
// any action that may initiate a Stripe payment transaction.
const cardAuthText = "  I authorize Acrobox to charge my card in accordance with the terms of service.\n"

const (
	colorINF = 34
	colorWRN = 33
	colorERR = 31
)

// Run parses the command line arguments, starting with the
// program name, and dispatches to the appropriate handler.
func Run(config *Config) error {
	c := &client{config: config}
	options := []cli.Option{
		cli.Scope("cli"),
		cli.Prefix("ACROBOX"),
		cli.Version(AppVersion),
		cli.Resolver(c.resolver),
		cli.Stdin(config.Stdin),
		cli.Stdout(config.Stdout),
		cli.Stderr(config.Stderr),
	}
	c.cli = cli.New(AppName, cli.NewUsageFS(docs.FS), []*cli.Flag{
		cli.NewFlag("host", &c.flags.host, cli.DefaultValue(username), cli.ShortFlag("h")),
		cli.NewFlag("verbose", &c.flags.verbose, cli.Bool(), cli.ShortFlag("v")),
		// Hidden
		cli.NewFlag("addr", &c.flags.addr, cli.DefaultValue("https://acrobox.io")),
		cli.NewFlag("port", &c.flags.port, cli.DefaultValue("22")),
	}, options...)
	c.cli.Use(func(next cli.Handler) cli.Handler {
		fn := func(args []string) error {
			c.http = newService(c.flags.addr, c.flags.auth)
			return next(args)
		}
		return cli.Handler(fn)
	})
	// Machine
	c.cli.Add("init", c.init, []*cli.Flag{
		cli.NewFlag("region", &c.flags.init.Region, cli.DefaultValue("nyc1"), cli.ShortFlag("r")),
		cli.NewFlag("size", &c.flags.init.Size, cli.DefaultValue("s-1vcpu-1gb-intel"), cli.ShortFlag("s")),
		cli.NewFlag("data-size", &c.flags.init.DataSize, cli.Kind(flagInt{}), cli.DefaultValue("1"), cli.ShortFlag("d")),
		cli.NewFlag("digitalocean-access-token", &c.flags.init.AccessToken, cli.EnvironmentKey("DIGITALOCEAN_ACCESS_TOKEN")),
		cli.NewFlag("token", &c.flags.auth),
		cli.NewFlag("force", &c.flags.init.force, cli.Bool(), cli.ShortFlag("f")),
	})
	c.cli.Add("cancel", c.cancel, []*cli.Flag{
		cli.NewFlag("token", &c.flags.auth),
		cli.NewFlag("force", &c.flags.cancel.force, cli.Bool(), cli.ShortFlag("f")),
	})
	c.cli.Add("renew", c.renew, []*cli.Flag{
		cli.NewFlag("token", &c.flags.auth),
		cli.NewFlag("force", &c.flags.renew.force, cli.Bool(), cli.ShortFlag("f")),
	})
	c.cli.Add("destroy", c.destroy, []*cli.Flag{
		cli.NewFlag("digitalocean-access-token", &c.flags.destroy.AccessToken, cli.EnvironmentKey("DIGITALOCEAN_ACCESS_TOKEN")),
		cli.NewFlag("token", &c.flags.auth),
		cli.NewFlag("force", &c.flags.destroy.force, cli.Bool(), cli.ShortFlag("f")),
	})
	c.cli.Add("ssh", c.ssh, nil)
	c.cli.Add("push", c.push, nil)
	c.cli.Add("pull", c.pull, nil)
	c.cli.Add("status", c.status, []*cli.Flag{
		cli.NewFlag("format", &c.flags.status.format, cli.DefaultValue("term"), cli.ShortFlag("f")),
	})
	c.cli.Add("metrics", c.metrics, []*cli.Flag{
		cli.NewFlag("format", &c.flags.metrics.format, cli.DefaultValue("term"), cli.ShortFlag("f")),
	})
	c.cli.Add("db/info", c.databaseInfo, nil)
	c.cli.Add("psql", c.psql, nil, cli.Proxy())
	c.cli.Add("redis-cli", c.redisCLI, nil, cli.Proxy())
	c.cli.Add("restore", c.restore, []*cli.Flag{
		cli.NewFlag("force", &c.flags.restore.force, cli.Bool(), cli.ShortFlag("f")),
	})
	commands := []string{
		"db/list",
		"db/create",
		"db/backup",
		"db/restore",
		"db/destroy",
		"backup",
		"update",
	}
	for _, cmd := range commands {
		c.proxy(cmd)
	}
	// Containers
	c.cli.Add("deploy", c.deploy, nil, cli.Proxy())
	c.cli.Add("logs", c.logs, nil, cli.Proxy())
	c.cli.Add("exec", c.exec, nil)
	commands = []string{
		"add",
		"remove",
		"list",
		"show",
		"run",
		"stop",
		"start",
		"reload",
		"restart",
		"env/set",
		"env/get",
		"env/all",
		"env/del",
	}
	for _, cmd := range commands {
		c.proxy(cmd)
	}
	return c.cli.Run(config.Args)
}

func (c *client) proxy(cmd string) *cli.Command {
	fn := func(args []string) error {
		return c.acroboxd(cmd, args)
	}
	return c.cli.Add(cmd, fn, nil, cli.Proxy())
}

func (c *client) init(args []string) error {
	if len(args) > 0 {
		return cli.ErrUsage
	}
	ipv4 := filepath.Join(c.config.Home, c.flags.host, "IPv4")
	_, err := os.Stat(ipv4)
	if !os.IsNotExist(err) {
		dir := filepath.Dir(ipv4)
		return fmt.Errorf("Machine '%s' already exists.", dir)
	}
	c.step(colorINF, "Creating a new key pair.")
	privateKey, authorizedKey, err := newKeyPair()
	if err != nil {
		c.step(colorERR, "%v", err)
		return cli.ErrExitFailure
	}
	c.flags.init.Product = "Indie Hacker"
	c.flags.init.Name = c.flags.host
	c.flags.init.PublicKey = string(authorizedKey)
	if !c.flags.init.force {
		c.step(colorWRN, "Payment authorization required.")
		c.cli.Printf(cardAuthText)
		name := c.cli.Prompt("\033[1;%dm•\033[0m \033[1;37mPlease type '%s' to agree:\033[0m ", colorWRN, c.flags.host)
		if name != c.flags.host {
			c.step(colorERR, "Input must be '%s' to agree.", c.flags.host)
			return cli.ErrExitFailure
		}
	}
	c.step(colorINF, "Initializing with '%s'.", c.flags.addr)
	id, err := c.http.initMachine(c.flags.init)
	if err != nil {
		verr, ok := err.(errorResponse)
		if ok {
			c.step(colorERR, verr.Message)
		} else {
			c.step(colorERR, "%v", err)
		}
		return cli.ErrExitFailure
	}
	c.step(colorINF, "Provisioning machine and associated resources.")
	m, err := c.waitForMachine(id)
	if err != nil {
		c.step(colorERR, "%v", err)
		return cli.ErrExitFailure
	}
	c.step(colorINF, "Writing machine configuration.")
	dir := filepath.Join(c.config.Home, c.flags.host)
	err = os.MkdirAll(dir, 0770)
	if err != nil {
		c.step(colorERR, "Data directory '%s' does not exist or cannot be created.", dir)
		return cli.ErrExitFailure
	}
	err = c.writeString("ID", m.ID)
	if err != nil {
		c.step(colorERR, "%v", err)
		return cli.ErrExitFailure
	}
	err = c.writeString("IPv4", m.IPv4)
	if err != nil {
		c.step(colorERR, "%v", err)
		return cli.ErrExitFailure
	}
	pub, _, _, _, err := ssh.ParseAuthorizedKey([]byte(m.PublicKey))
	if err != nil {
		c.step(colorERR, "%v", err)
		return cli.ErrExitFailure
	}
	knownHosts := knownhosts.Line([]string{m.IPv4 + ":" + c.flags.port}, pub)
	err = c.writeString("known_hosts", knownHosts)
	if err != nil {
		c.step(colorERR, "%v", err)
		return cli.ErrExitFailure
	}
	err = c.writeKey("id_ed25519", privateKey)
	if err != nil {
		c.step(colorERR, "%v", err)
		return cli.ErrExitFailure
	}
	err = c.writeKey("id_ed25519.pub", authorizedKey)
	if err != nil {
		c.step(colorERR, "%v", err)
		return cli.ErrExitFailure
	}
	c.step(colorINF, "Waiting for SSH connectivity.")
	err = c.waitForSSH(m.IPv4)
	if err != nil {
		c.step(colorERR, "%v", err)
		return cli.ErrExitFailure
	}
	c.step(colorINF, "Waiting for machine setup.")
	err = c.waitForAcrobox()
	if err != nil {
		c.step(colorERR, "%v", err)
		return cli.ErrExitFailure
	}
	c.step(colorINF, "Waiting for service setup.")
	err = c.waitForService()
	if err != nil {
		c.step(colorERR, "%v", err)
		return cli.ErrExitFailure
	}
	c.step(colorINF, "Acrobox is ready.")
	return nil
}

func (c *client) cancel(args []string) error {
	if len(args) > 0 {
		return cli.ErrUsage
	}
	if !c.flags.cancel.force {
		c.cli.Printf("Confirmation to cancel service '%s' is required.\n", c.flags.host)
		err := c.promptToAgree()
		if err != nil {
			return err
		}
	}
	id, err := c.getID()
	if err != nil {
		return err
	}
	return c.http.cancelMachine(id)
}

func (c *client) renew(args []string) error {
	if len(args) > 0 {
		return cli.ErrUsage
	}
	if !c.flags.renew.force {
		c.cli.Printf("Payment authorization required.\n")
		c.cli.Printf(cardAuthText)
		err := c.promptToAgree()
		if err != nil {
			return err
		}
	}
	id, err := c.getID()
	if err != nil {
		return err
	}
	return c.http.renewMachine(id)
}

func (c *client) destroy(args []string) error {
	if len(args) > 0 {
		return cli.ErrUsage
	}
	if !c.flags.destroy.force {
		c.cli.Printf("Confirmation to destroy machine '%s' is required.\n", c.flags.host)
		c.cli.Printf("  All data will be lost.\n")
		c.cli.Printf("  This action cannot be undone.\n")
		err := c.promptToAgree()
		if err != nil {
			return err
		}
	}
	id, err := c.getID()
	if err != nil {
		return err
	}
	err = c.http.destroyMachine(id, c.flags.destroy)
	if err != nil {
		return err
	}
	dir := filepath.Join(c.config.Home, c.flags.host)
	return os.RemoveAll(dir)
}

func (c *client) ssh(args []string) error {
	if len(args) > 0 {
		return c.runWithOutput(args[0], args[1:]...)
	}
	session, err := c.newSession()
	if err != nil {
		return err
	}
	defer session.Close()
	session.Stdin = c.config.Stdin
	session.Stdout = c.config.Stdout
	session.Stderr = c.config.Stderr
	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
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
	err = session.Shell()
	if err != nil {
		return err
	}
	err = session.Wait()
	if err != nil {
		_, ok := err.(*ssh.ExitError)
		if ok {
			return cli.ErrExitFailure
		}
		return err
	}
	return nil
}

func (c *client) push(args []string) error {
	if len(args) < 2 {
		return cli.ErrUsage
	}
	target := args[len(args)-1]
	source := args[:len(args)-1]
	fileType, _, err := c.run("stat", "-L", "-c", "%F", target)
	if err != nil {
		serr, ok := err.(*ssh.ExitError)
		if !ok || serr.ExitStatus() != 1 {
			return err
		}
	}
	targetIsDir := string(fileType) == "directory\n"
	if targetIsDir {
		for _, s := range source {
			err = c.pushOne(s, target, true)
			if err != nil {
				return err
			}
		}
		return nil
	}
	if len(source) > 1 {
		return cli.ErrUsage
	}
	return c.pushOne(source[0], target, false)
}

func (c *client) pushOne(source, target string, targetIsDir bool) error {
	f, err := os.Open(source)
	if err != nil {
		return err
	}
	defer f.Close()
	fi, err := f.Stat()
	if err != nil {
		return err
	}
	if fi.IsDir() {
		source = strings.TrimSuffix(source, "/")
		return c.pushDir(source, target)
	}
	filename := target
	if targetIsDir {
		filename = filepath.Join(target, fi.Name())
	}
	_, _, err = c.runWithStdin(f, "sh", "-c", "cat > "+filename)
	return err
}

func (c *client) pushDir(source, target string) error {
	dir := filepath.Dir(source)
	fn := func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		name := strings.TrimPrefix(path, dir)
		filename := filepath.Join(target, name)
		if fi.IsDir() {
			if filename == target {
				return nil
			}
			_, _, err = c.run("mkdir", filename)
			return err
		}
		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()
		_, _, err = c.runWithStdin(f, "sh", "-c", "cat > "+filename)
		return err
	}
	return filepath.Walk(source, fn)
}

func (c *client) pull(args []string) error {
	if len(args) < 2 {
		return cli.ErrUsage
	}
	target, err := filepath.Abs(args[len(args)-1])
	if err != nil {
		return err
	}
	source := args[:len(args)-1]
	fi, err := os.Stat(target)
	if err == nil && fi.IsDir() {
		for _, s := range source {
			err = c.pullOne(s, target)
			if err != nil {
				return err
			}
		}
		return nil
	}
	if len(source) > 1 {
		return cli.ErrUsage
	}
	return c.pullOne(source[0], target)
}

func (c *client) pullOne(source, target string) error {
	fileType, _, err := c.run("stat", "-L", "-c", "%F", source)
	if err != nil {
		serr, ok := err.(*ssh.ExitError)
		if !ok || serr.ExitStatus() != 1 {
			return err
		}
	}
	target = filepath.Join(target, filepath.Base(source))
	sourceIsDir := string(fileType) == "directory\n"
	if sourceIsDir {
		source = strings.TrimSuffix(source, "/")
		return c.pullDir(source, target)
	}
	b, _, err := c.run("cat", source)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(target, b, 0640)
}

func (c *client) pullDir(source, target string) error {
	b, _, err := c.run("find", source)
	if err != nil {
		return err
	}
	err = os.MkdirAll(target, 0750)
	if err != nil {
		return err
	}
	s := bufio.NewScanner(bytes.NewReader(b))
	for s.Scan() {
		filename := s.Text()
		if filename == source {
			continue
		}
		err = c.pullOne(filename, target)
		if err != nil {
			return err
		}
	}
	return s.Err()
}

func (c *client) status(args []string) error {
	if len(args) > 0 {
		return cli.ErrUsage
	}
	ipv4, err := c.getIPv4()
	if err != nil {
		return err
	}
	stdout, stderr, err := c.run("docker", "exec", "acroboxd", "acroboxd", "status")
	if err != nil {
		c.cli.Errorf("%s\n", stderr)
		return cli.ErrExitFailure
	}
	data := acroboxdStatus{}
	err = json.Unmarshal(stdout, &data)
	if err != nil {
		return err
	}
	switch c.flags.status.format {
	case "term":
	case "json":
		return json.NewEncoder(c.config.Stdout).Encode(data)
	default:
		return fmt.Errorf("Format must be 'term' or 'json'.")
	}
	now := time.Now()
	c.cli.Printf("\n")
	c.cli.Printf("System\n")
	c.cli.Printf("  IPv4:       %s\n", ipv4)
	c.cli.Printf("  Hostname:   %s\n", c.flags.host)
	c.cli.Printf("  Booted At:  %s%s\n", formatTime(data.SystemBootedAt), formatDuration(data.SystemBootedAt, now))
	c.cli.Printf("  Started At: %s%s\n", formatTime(data.AcroboxStartedAt), formatDuration(data.AcroboxStartedAt, now))
	c.cli.Printf("\n")
	c.cli.Printf("Memory\n")
	c.cli.Printf("  Go Routines: %d\n", data.MemStatsGoRoutines)
	c.cli.Printf("  Total: %d bytes\n", data.MemStatsMemoryTotal)
	c.cli.Printf("  Alloc: %d bytes\n", data.MemStatsMemoryAlloc)
	c.cli.Printf("  Count: %d\n", data.MemStatsMemoryCount)
	c.cli.Printf("\n")
	c.cli.Printf("Requests\n")
	c.cli.Printf("  Rate:   %d req/s", data.ProxyRequestsPerSec)
	if data.ProxyRequestsPerSecPeak > 0 {
		c.cli.Printf(" (peak %d at %v)\n", data.ProxyRequestsPerSecPeak, formatTime(data.ProxyRequestsPerSecPeakAt))
	} else {
		c.cli.Printf("\n")
	}
	c.cli.Printf("  Active: %d", data.ProxyActiveRequests)
	if data.ProxyActiveRequestsPeak > 0 {
		c.cli.Printf(" (peak %d at %v)\n", data.ProxyActiveRequestsPeak, formatTime(data.ProxyActiveRequestsPeakAt))
	} else {
		c.cli.Printf("\n")
	}
	c.cli.Printf("\n")
	c.cli.Printf("Latency p99 (window=10m interval=10s)\n")
	c.formatLatency(data.ProxyLatencyP99, 10)
	c.cli.Printf("\n")
	c.cli.Printf("Backups\n")
	c.cli.Printf("  Next: %s%s\n", formatTime(data.NextBackupAt), formatDuration(data.NextBackupAt, now))
	if !data.LastBackupAt.IsZero() {
		c.cli.Printf("  Last: %s (took %s)\n", formatTime(data.LastBackupAt), data.LastBackupTook)
	}
	c.cli.Printf("\n")
	c.cli.Printf("Updates\n")
	c.cli.Printf("  Next: %s%s\n", formatTime(data.NextUpdateAt), formatDuration(data.NextUpdateAt, now))
	if !data.LastUpdateAt.IsZero() {
		c.cli.Printf("  Last: %s (took %s)\n", formatTime(data.LastUpdateAt), data.LastUpdateTook)
	}
	c.cli.Printf("\n")
	if !data.LicenseUpdatedAt.IsZero() {
		c.cli.Printf("License\n")
		c.cli.Printf("  Updated At: %s%s\n", formatTime(data.LicenseUpdatedAt), formatDuration(data.LicenseUpdatedAt, now))
		c.cli.Printf("\n")
	}
	return nil
}

func (c *client) formatLatency(vs []float64, height int) {
	max := 0
	for _, n := range vs {
		if int(n) > max {
			max = int(n)
		}
	}
	ys := make([]int, len(vs))
	if max != 0 {
		for i, v := range vs {
			ys[i] = (int(v)*height + max/2) / max
		}
	}
	for h := height; h > 0; h-- {
		s := ""
		for i := range vs {
			if ys[i] >= h {
				s += "█"
			} else {
				s += " "
			}
		}
		if h == height {
			c.cli.Printf("  %s %d ms\n", s, max)
		} else {
			c.cli.Printf("  %s\n", s)
		}
	}
}

func (c *client) metrics(args []string) error {
	if len(args) > 0 {
		return cli.ErrUsage
	}
	stdout, stderr, err := c.run("docker", "exec", "acroboxd", "acroboxd", "metrics")
	if err != nil {
		c.cli.Errorf("%s\n", stderr)
		return cli.ErrExitFailure
	}
	data := make(map[string]interface{})
	err = json.Unmarshal(stdout, &data)
	if err != nil {
		return err
	}
	switch c.flags.metrics.format {
	case "term":
		return fmt.Errorf("Format 'term' not implemented.")
	case "json":
		return json.NewEncoder(c.config.Stdout).Encode(data)
	default:
		return fmt.Errorf("Format must be 'term' or 'json'.")
	}
	return nil
}

func (c *client) databaseInfo(args []string) error {
	if len(args) > 0 {
		return cli.ErrUsage
	}
	ipv4, err := c.getIPv4()
	if err != nil {
		return err
	}
	password, _, err := c.run(`jq -r '.environment."acrobox/acroboxd".POSTGRES_PASSWORD' /acrobox/config.json`)
	if err != nil {
		return err
	}
	keyFile := filepath.Join(c.config.Home, c.flags.host, "id_ed25519")
	c.cli.Printf("Host           %s\n", ipv4)
	c.cli.Printf("Username       %s\n", username)
	c.cli.Printf("Password       %s\n", strings.TrimSpace(string(password)))
	c.cli.Printf("SSH Key File   %s\n", keyFile)
	return nil
}

func (c *client) psql(args []string) error {
	args = append([]string{"exec", "-i", "-t", "-u", "postgres", "postgres", "psql"}, args...)
	return c.runWithOutput("docker", args...)
}

func (c *client) redisCLI(args []string) error {
	args = append([]string{"exec", "-i", "-t", "-u", "redis", "redis", "redis-cli"}, args...)
	return c.runWithOutput("docker", args...)
}

func (c *client) restore(args []string) error {
	if !c.flags.restore.force {
		c.cli.Printf("Confirmation to restore machine '%s' is required.\n", c.flags.host)
		c.cli.Printf("  This action has the potential to overwrite data with old data.\n")
		c.cli.Printf("  This action cannot be undone.\n")
		err := c.promptToAgree()
		if err != nil {
			return err
		}
	}
	args = append([]string{"exec", "acroboxd", "acroboxd", "restore"}, args...)
	stdout, stderr, err := c.run("docker", args...)
	if err != nil {
		c.cli.Errorf("%s\n", stderr)
		return cli.ErrExitFailure
	}
	data := make(map[string]string)
	err = json.Unmarshal(stdout, &data)
	if err != nil {
		return err
	}
	containerID := strings.TrimSpace(string(data["container_id"]))
	return c.runWithOutput("docker", "logs", "-f", containerID)
}

func (c *client) deploy(args []string) error {
	if len(args) < 1 {
		return cli.ErrUsage
	}
	image := imageName(args[len(args)-1])
	f, err := os.CreateTemp("", "abx-deploy-")
	if err != nil {
		return err
	}
	local := f.Name()
	defer os.Remove(local)
	cmd := exec.Command("docker", "save", image, "-o", local)
	cmd.Stdout = c.config.Stdout
	cmd.Stderr = c.config.Stderr
	err = cmd.Run()
	if err != nil {
		return err
	}
	err = c.pushOne(local, "/tmp/", true)
	if err != nil {
		return err
	}
	base := filepath.Base(local)
	remote := filepath.Join("/tmp", base)
	_, stderr, err := c.run("sh", "-c", "docker load -i "+remote+" && rm "+remote)
	if err != nil {
		c.cli.Errorf("%s\n", stderr)
		return err
	}
	return c.acroboxd("deploy", args)
}

func (c *client) logs(args []string) error {
	args = append([]string{"logs"}, args...)
	return c.runWithOutput("docker", args...)
}

func (c *client) exec(args []string) error {
	if len(args) < 2 {
		return cli.ErrUsage
	}
	args = append([]string{"exec", "-i", "-t", args[0], args[1]}, args[2:]...)
	return c.runWithOutput("docker", args...)
}

func (c *client) step(level int, format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	c.cli.Printf("\033[1;%dm•\033[0m \033[1;37m%s\033[0m\n", level, message)
}

func (c *client) promptToAgree() error {
	name := c.cli.Prompt("Please type '%s' to agree: ", c.flags.host)
	if name != c.flags.host {
		return fmt.Errorf("Input must be '%s' to agree.", c.flags.host)
	}
	return nil
}

func (c *client) resolver(err error) {
	switch e := err.(type) {
	case *fs.PathError:
		prefix := filepath.Join(c.config.Home, c.flags.host)
		if strings.HasPrefix(e.Path, prefix) {
			c.cli.Errorf("Machine '%s' does not exist.\n", prefix)
		}
	case *ssh.ExitError:
		// no-op
	default:
		c.cli.Errorf("%v\n", err)
	}
}

func formatTime(t time.Time) string {
	return t.Format("Mon 2006-01-02 15:04:05 MST")
}

func formatDuration(t, now time.Time) string {
	if t.IsZero() || t.Equal(now) {
		return ""
	}
	if t.After(now) {
		d := t.Sub(now)
		return fmt.Sprintf(" (%s from now)", roundDuration(d))
	}
	d := now.Sub(t)
	return fmt.Sprintf(" (%s ago)", roundDuration(d))
}

func roundDuration(d time.Duration) time.Duration {
	if d < time.Second {
		return d.Round(time.Millisecond)
	}
	return d.Round(time.Second)
}

func imageName(s string) string {
	for _, c := range []string{":", "@"} {
		i := strings.Index(s, c)
		if i != -1 {
			s = s[:i]
		}
	}
	return s
}
