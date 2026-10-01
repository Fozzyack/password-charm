// Package githubsync backs up and restores encrypted stores using GitHub Git repositories.
package githubsync

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/ProtonMail/gopenpgp/v3/crypto"
)

const snapshotDir = "passwords"

var ErrNotConfigured = errors.New("configure a GitHub repository first")

// Client uses existing Git SSH keys or credential helpers; it stores no credentials.
type Client struct {
	remote string
}

var repositoryPath = regexp.MustCompile(`^[A-Za-z0-9_.-]+/[A-Za-z0-9_.-]+$`)

// New accepts credential-free GitHub SSH or HTTPS repository URLs.
func New(remote string) (*Client, error) {
	remote = strings.TrimSpace(remote)
	var path string
	if strings.HasPrefix(remote, "git@github.com:") {
		path = strings.TrimPrefix(remote, "git@github.com:")
	} else {
		u, err := url.Parse(remote)
		if err != nil || u.Host != "github.com" || u.RawQuery != "" || u.Fragment != "" {
			return nil, errors.New("use a GitHub SSH or HTTPS repository URL")
		}
		switch u.Scheme {
		case "https":
			if u.User != nil {
				return nil, errors.New("use Git credentials instead of embedding credentials in the URL")
			}
		case "ssh":
			if u.User == nil || u.User.String() != "git" {
				return nil, errors.New("GitHub SSH URLs must use the git user")
			}
		default:
			return nil, errors.New("use a GitHub SSH or HTTPS repository URL")
		}
		path = strings.TrimPrefix(u.Path, "/")
	}
	path = strings.TrimSuffix(path, ".git")
	if !repositoryPath.MatchString(path) {
		return nil, errors.New("repository URL must identify github.com/owner/repository")
	}
	for _, part := range strings.Split(path, "/") {
		if part == "." || part == ".." {
			return nil, errors.New("invalid repository path")
		}
	}
	return &Client{remote: remote}, nil
}

// SaveConfig saves the repository URL outside the password store.
func SaveConfig(path, remote string) error {
	c, err := New(remote)
	if err != nil {
		return err
	}
	data, err := json.Marshal(struct {
		Repository string `json:"repository"`
	}{c.remote})
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}
	f, err := os.CreateTemp(filepath.Dir(path), ".github-config-*")
	if err != nil {
		return err
	}
	defer os.Remove(f.Name())
	if _, err := f.Write(data); err != nil {
		f.Close()
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	return os.Rename(f.Name(), path)
}

func LoadConfig(path string) (*Client, error) {
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil, ErrNotConfigured
	}
	if err != nil {
		return nil, err
	}
	var config struct {
		Repository string `json:"repository"`
	}
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("read GitHub configuration: %w", err)
	}
	return New(config.Repository)
}

// Backup publishes a complete snapshot, including deletions, on the remote default branch.
// Only encrypted .gpg files and .checker/init.gpg enter the snapshot.
func (c *Client) Backup(ctx context.Context, store string) error {
	files, err := encryptedFiles(store)
	if err != nil {
		return err
	}
	repo, cleanup, err := c.clone(ctx)
	if err != nil {
		return err
	}
	defer cleanup()
	if _, err := git(ctx, repo, "rev-parse", "--verify", "HEAD"); err != nil {
		if _, err := git(ctx, repo, "symbolic-ref", "HEAD", "refs/heads/main"); err != nil {
			return err
		}
	}
	destination := filepath.Join(repo, snapshotDir)
	if err := os.RemoveAll(destination); err != nil {
		return err
	}
	if err := writeSnapshot(destination, files); err != nil {
		return err
	}
	if _, err := git(ctx, repo, "add", "--all", "--", snapshotDir); err != nil {
		return err
	}
	changes, err := git(ctx, repo, "diff", "--cached", "--name-only")
	if err != nil {
		return err
	}
	if strings.TrimSpace(changes) == "" {
		return nil
	}
	if _, err := git(ctx, repo, "-c", "user.name=Password Manager", "-c", "user.email=password-manager@users.noreply.github.com", "-c", "commit.gpgsign=false", "commit", "-m", "Back up encrypted password store"); err != nil {
		return err
	}
	_, err = git(ctx, repo, "push", "origin", "HEAD")
	return err
}

// Restore installs a snapshot into an empty store. Existing passwords are never overwritten.
// After restoring, log in with the backed-up store's master password.
func (c *Client) Restore(ctx context.Context, store string) error {
	if err := emptyStore(store); err != nil {
		return err
	}
	repo, cleanup, err := c.clone(ctx)
	if err != nil {
		return err
	}
	defer cleanup()
	files, err := encryptedFiles(filepath.Join(repo, snapshotDir))
	if err != nil {
		return fmt.Errorf("read GitHub backup: %w", err)
	}
	// Stage beside the destination so renaming remains on the same filesystem.
	stage, err := os.MkdirTemp(filepath.Dir(store), ".password-manager-restore-*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(stage)
	if err := writeSnapshot(stage, files); err != nil {
		return err
	}
	if err := emptyStore(store); err != nil {
		return err
	}
	old := stage + "-previous"
	if err := os.Rename(store, old); err != nil {
		return err
	}
	if err := os.Rename(stage, store); err != nil {
		if rollbackErr := os.Rename(old, store); rollbackErr != nil {
			return fmt.Errorf("restore failed: %v; original store remains at %s: %w", err, old, rollbackErr)
		}
		return err
	}
	return os.RemoveAll(old)
}

func (c *Client) clone(ctx context.Context) (string, func(), error) {
	root, err := os.MkdirTemp("", "password-manager-github-*")
	if err != nil {
		return "", nil, err
	}
	cleanup := func() { os.RemoveAll(root) }
	repo := filepath.Join(root, "repository")
	if _, err := git(ctx, "", "clone", "--", c.remote, repo); err != nil {
		cleanup()
		return "", nil, err
	}
	return repo, cleanup, nil
}

func git(ctx context.Context, dir string, args ...string) (string, error) {
	args = append([]string{"-c", "core.hooksPath=/dev/null"}, args...)
	cmd := exec.CommandContext(ctx, "git", args...)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "GIT_TERMINAL_PROMPT=0")
	output, err := cmd.CombinedOutput()
	if err != nil {
		if ctx.Err() != nil {
			return "", fmt.Errorf("GitHub operation: %w", ctx.Err())
		}
		return "", fmt.Errorf("git: %w: %s", err, strings.TrimSpace(string(output)))
	}
	return string(output), nil
}

func encryptedFiles(root string) (map[string][]byte, error) {
	files := make(map[string][]byte)
	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		if entry.Type()&os.ModeSymlink != 0 {
			return fmt.Errorf("symlinks are not supported in password snapshots: %s", rel)
		}
		if entry.IsDir() {
			if rel == "." || rel == ".checker" {
				return nil
			}
			return filepath.SkipDir
		}
		if rel != filepath.Join(".checker", "init.gpg") && (filepath.Dir(rel) != "." || !strings.HasSuffix(rel, ".gpg")) {
			return nil
		}
		if !entry.Type().IsRegular() {
			return fmt.Errorf("not a regular encrypted file: %s", rel)
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		message, err := crypto.NewPGPMessageFromArmored(string(data))
		if err != nil || len(message.KeyPacket) == 0 || len(message.DataPacket) == 0 {
			return fmt.Errorf("file is not an encrypted OpenPGP message: %s", rel)
		}
		files[rel] = data
		return nil
	})
	if err != nil {
		return nil, err
	}
	if _, ok := files[filepath.Join(".checker", "init.gpg")]; !ok {
		return nil, errors.New("snapshot is missing the encrypted login checker (.checker/init.gpg)")
	}
	return files, nil
}

func writeSnapshot(root string, files map[string][]byte) error {
	if err := os.MkdirAll(filepath.Join(root, ".checker"), 0700); err != nil {
		return err
	}
	for name, data := range files {
		if err := os.WriteFile(filepath.Join(root, name), data, 0600); err != nil {
			return err
		}
	}
	return nil
}

func emptyStore(root string) error {
	return filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		if entry.IsDir() && (rel == "." || rel == ".checker") {
			return nil
		}
		return errors.New("restore requires an empty password store; existing files will not be overwritten")
	})
}
