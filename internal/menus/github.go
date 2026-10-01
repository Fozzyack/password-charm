package menus

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/Fozzyack/password-manager/internal/githubsync"
	"github.com/Fozzyack/password-manager/internal/ui"
	uimenu "github.com/Fozzyack/password-manager/internal/ui/menu"
	"github.com/Fozzyack/password-manager/internal/ui/textinput"
)

func githubConfigPath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "password-manager", "github.json"), nil
}

// ConfigureGitHub asks for an existing GitHub repository URL.
func (m *Menu) ConfigureGitHub() (bool, error) {
	path, err := githubConfigPath()
	if err != nil {
		return false, err
	}
	remote := ""
	m.Options.ErrorMessage = ""
	_, err = ui.Run(textinput.InitialModelWithMasking(
		"GitHub backup: enter your private repository URL",
		"git@github.com:owner/repository.git", &remote, m.Options, false))
	if err != nil || m.Options.Quit {
		return false, err
	}
	if err := githubsync.SaveConfig(path, remote); err != nil {
		return false, err
	}
	return true, nil
}

func (m *Menu) githubClient() (*githubsync.Client, error) {
	path, err := githubConfigPath()
	if err != nil {
		return nil, err
	}
	client, err := githubsync.LoadConfig(path)
	if errors.Is(err, githubsync.ErrNotConfigured) {
		saved, configureErr := m.ConfigureGitHub()
		if configureErr != nil || !saved {
			return nil, configureErr
		}
		return githubsync.LoadConfig(path)
	}
	return client, err
}

// BackupToGitHub publishes the current encrypted store on demand.
func (m *Menu) BackupToGitHub() (bool, error) {
	client, err := m.githubClient()
	if err != nil || client == nil {
		return false, err
	}
	fmt.Println("Saving encrypted passwords to GitHub...")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	if err := client.Backup(ctx, m.passwordFolder.FolderLocation); err != nil {
		return false, err
	}
	return true, nil
}

// PrepareStore offers restore before a new machine creates its login checker.
func (m *Menu) PrepareStore() error {
	if m.passwordFolder.InitCheck {
		return nil
	}
	for {
		final, err := ui.Run(uimenu.InitialSetupMenuModel(m.Options))
		if err != nil || m.Options.Quit {
			return err
		}
		action := final.(uimenu.MenuModel).GetSelectedAction()
		if action == "github_configure" {
			if _, err := m.ConfigureGitHub(); err != nil || m.Options.Quit {
				return err
			}
			continue
		}
		if action != "github_restore" {
			return nil
		}
		break
	}
	client, err := m.githubClient()
	if err != nil || client == nil {
		return err
	}
	fmt.Println("Restoring encrypted passwords from GitHub...")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	if err := client.Restore(ctx, m.passwordFolder.FolderLocation); err != nil {
		return err
	}
	m.passwordFolder.InitCheck = true
	return m.passwordFolder.RefreshDirectoryListing()
}
