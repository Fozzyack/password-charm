package githubsync

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/Fozzyack/password-manager/internal/encryption"
	"github.com/Fozzyack/password-manager/internal/fileio"
)

func TestRepositoryURLs(t *testing.T) {
	for _, remote := range []string{"git@github.com:owner/vault.git", "https://github.com/owner/vault", "ssh://git@github.com/owner/vault.git"} {
		if _, err := New(remote); err != nil {
			t.Errorf("New(%q): %v", remote, err)
		}
	}
	for _, remote := range []string{"", "-bad", "https://example.com/owner/vault", "https://github.com.evil.test/owner/vault", "https://token@github.com/owner/vault", "https://github.com/owner/vault?token=secret", "git@github.com:../vault", "file:///tmp/repo", "ssh://user@github.com/owner/vault"} {
		if _, err := New(remote); err == nil {
			t.Errorf("New(%q) accepted invalid URL", remote)
		}
	}
}

func TestConfig(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config", "github.json")
	if _, err := LoadConfig(path); !errors.Is(err, ErrNotConfigured) {
		t.Fatalf("missing configuration: %v", err)
	}
	remote := "git@github.com:owner/vault.git"
	if err := SaveConfig(path, remote); err != nil {
		t.Fatal(err)
	}
	client, err := LoadConfig(path)
	if err != nil || client.remote != remote {
		t.Fatalf("configuration round trip: %v, %v", client, err)
	}
	info, err := os.Stat(path)
	if err != nil || info.Mode().Perm() != 0600 {
		t.Fatalf("configuration permissions: %v, %v", info, err)
	}
	if err := SaveConfig(path, "https://token@github.com/owner/vault"); err == nil {
		t.Fatal("saved embedded credentials")
	}
}

func TestBackupAndRestore(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git is not installed")
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	root := t.TempDir()
	remote := filepath.Join(root, "remote.git")
	if _, err := git(ctx, "", "init", "--bare", "--initial-branch=main", remote); err != nil {
		t.Fatal(err)
	}
	// A local bare repository exercises Git without network access or credentials.
	client := &Client{remote: remote}
	store := filepath.Join(root, "store")
	makeStore(t, store)
	pf := &fileio.PasswordFolder{FolderLocation: store, Password: "test master password"}
	ef := encryption.NewEncryption(pf)
	for _, name := range []string{".checker/init", "example", "deleted"} {
		if err := ef.EncryptPasswordAndWriteToFile(name, encryption.Data{Password: "encrypted test value"}); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.WriteFile(filepath.Join(store, "notes.txt"), []byte("must not be uploaded"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := client.Backup(ctx, store); err != nil {
		t.Fatal(err)
	}
	firstHead, err := git(ctx, remote, "rev-parse", "HEAD")
	if err != nil {
		t.Fatal(err)
	}
	if err := client.Backup(ctx, store); err != nil {
		t.Fatal(err)
	}
	secondHead, _ := git(ctx, remote, "rev-parse", "HEAD")
	if firstHead != secondHead {
		t.Fatal("unchanged backup created a commit")
	}
	// Preserve repository content outside the managed snapshot.
	clone, cleanup, err := client.clone(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer cleanup()
	if err := os.WriteFile(filepath.Join(clone, "README.md"), []byte("repository notes"), 0600); err != nil {
		t.Fatal(err)
	}
	for _, args := range [][]string{{"add", "README.md"}, {"-c", "user.name=Test", "-c", "user.email=test@example.com", "-c", "commit.gpgsign=false", "commit", "-m", "Add notes"}, {"push", "origin", "HEAD"}} {
		if _, err := git(ctx, clone, args...); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.Remove(filepath.Join(store, "deleted.gpg")); err != nil {
		t.Fatal(err)
	}
	if err := client.Backup(ctx, store); err != nil {
		t.Fatal(err)
	}
	tree, err := git(ctx, remote, "ls-tree", "-r", "--name-only", "HEAD")
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(tree, "notes.txt") || strings.Contains(tree, "deleted.gpg") || !strings.Contains(tree, "README.md") || !strings.Contains(tree, "passwords/.checker/init.gpg") {
		t.Fatalf("unexpected snapshot contents: %s", tree)
	}
	restored := filepath.Join(root, "restored")
	makeStore(t, restored)
	if err := client.Restore(ctx, restored); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{".checker/init.gpg", "example.gpg"} {
		original, _ := os.ReadFile(filepath.Join(store, name))
		copy, err := os.ReadFile(filepath.Join(restored, name))
		if err != nil || string(original) != string(copy) {
			t.Fatalf("restore changed %s: %v", name, err)
		}
		info, _ := os.Stat(filepath.Join(restored, name))
		if info.Mode().Perm() != 0600 {
			t.Fatalf("restored file permissions: %v", info.Mode())
		}
	}
	restoredEncryption := encryption.NewEncryption(&fileio.PasswordFolder{FolderLocation: restored, Password: pf.Password})
	if data, err := restoredEncryption.DecryptPasswordFromFile("example"); err != nil || data.Password != "encrypted test value" {
		t.Fatalf("restored password cannot be decrypted: %v", err)
	}
	if err := client.Restore(ctx, restored); err == nil {
		t.Fatal("restore overwrote an initialized store")
	}
}

func TestInvalidSnapshots(t *testing.T) {
	root := t.TempDir()
	makeStore(t, root)
	if _, err := encryptedFiles(root); err == nil {
		t.Fatal("accepted missing checker")
	}
	if err := os.WriteFile(filepath.Join(root, ".checker", "init.gpg"), []byte("plaintext"), 0600); err != nil {
		t.Fatal(err)
	}
	if _, err := encryptedFiles(root); err == nil {
		t.Fatal("accepted plaintext named .gpg")
	}
	if err := os.Remove(filepath.Join(root, ".checker", "init.gpg")); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(filepath.Join(root, "outside"), filepath.Join(root, "secret.gpg")); err != nil {
		t.Fatal(err)
	}
	if _, err := encryptedFiles(root); err == nil {
		t.Fatal("accepted a symlink")
	}
	if err := emptyStore(root); err == nil {
		t.Fatal("accepted nonempty destination")
	}
}

func TestFailedRestoreLeavesStoreIntact(t *testing.T) {
	root := t.TempDir()
	store := filepath.Join(root, "store")
	makeStore(t, store)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	client := &Client{remote: filepath.Join(root, "missing.git")}
	if err := client.Restore(ctx, store); err == nil {
		t.Fatal("restore succeeded with unavailable remote")
	}
	if err := emptyStore(store); err != nil {
		t.Fatalf("failed restore changed the destination: %v", err)
	}
}

func makeStore(t *testing.T, path string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Join(path, ".checker"), 0700); err != nil {
		t.Fatal(err)
	}
}
