package ssh

import (
	"fmt"
	"net"
	"os"
	"time"

	"golang.org/x/crypto/ssh"
)

// Client wraps an SSH connection to a remote server.
type Client struct {
	conn   *ssh.Client
	Host   string
}

// Config holds the parameters needed to establish an SSH connection.
type Config struct {
	Host       string
	Port       int
	User       string
	KeyPath    string
	Timeout    time.Duration
}

// New establishes a new SSH connection using the provided config.
func New(cfg Config) (*Client, error) {
	if cfg.Port == 0 {
		cfg.Port = 22
	}
	if cfg.Timeout == 0 {
		cfg.Timeout = 15 * time.Second
	}

	key, err := os.ReadFile(cfg.KeyPath)
	if err != nil {
		return nil, fmt.Errorf("reading key %q: %w", cfg.KeyPath, err)
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, fmt.Errorf("parsing private key: %w", err)
	}

	sshCfg := &ssh.ClientConfig{
		User:            cfg.User,
		Auth:            []ssh.AuthMethod{ssh.PublicKeys(signer)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // TODO: replace with known_hosts
		Timeout:         cfg.Timeout,
	}

	addr := net.JoinHostPort(cfg.Host, fmt.Sprintf("%d", cfg.Port))
	conn, err := ssh.Dial("tcp", addr, sshCfg)
	if err != nil {
		return nil, fmt.Errorf("dialing %s: %w", addr, err)
	}

	return &Client{conn: conn, Host: cfg.Host}, nil
}

// RunCommand executes a shell command on the remote server and returns combined output.
func (c *Client) RunCommand(cmd string) (string, error) {
	sess, err := c.conn.NewSession()
	if err != nil {
		return "", fmt.Errorf("creating session: %w", err)
	}
	defer sess.Close()

	out, err := sess.CombinedOutput(cmd)
	if err != nil {
		return string(out), fmt.Errorf("running command %q: %w", cmd, err)
	}
	return string(out), nil
}

// Close terminates the underlying SSH connection.
func (c *Client) Close() error {
	return c.conn.Close()
}
