package ssh

import (
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/pkg/sftp"
)

// UploadFile copies the content of src to remotePath on the server.
func (c *Client) UploadFile(remotePath string, src io.Reader) error {
	client, err := sftp.NewClient(c.conn)
	if err != nil {
		return fmt.Errorf("creating sftp client: %w", err)
	}
	defer client.Close()

	dir := filepath.Dir(remotePath)
	if err := client.MkdirAll(dir); err != nil {
		return fmt.Errorf("creating remote dir %q: %w", dir, err)
	}

	dst, err := client.Create(remotePath)
	if err != nil {
		return fmt.Errorf("creating remote file %q: %w", remotePath, err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return fmt.Errorf("uploading to %q: %w", remotePath, err)
	}
	return nil
}

// UploadString writes content as a string to remotePath on the server.
func (c *Client) UploadString(remotePath, content string) error {
	return c.UploadFile(remotePath, strings.NewReader(content))
}

// DownloadFile retrieves a remote file and writes its contents to dst.
func (c *Client) DownloadFile(remotePath string, dst io.Writer) error {
	client, err := sftp.NewClient(c.conn)
	if err != nil {
		return fmt.Errorf("creating sftp client: %w", err)
	}
	defer client.Close()

	src, err := client.Open(remotePath)
	if err != nil {
		return fmt.Errorf("opening remote file %q: %w", remotePath, err)
	}
	defer src.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return fmt.Errorf("downloading %q: %w", remotePath, err)
	}
	return nil
}
