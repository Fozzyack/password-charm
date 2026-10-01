package menu

import "github.com/Fozzyack/password-manager/internal/types"

// InitialSetupMenuModel lets a new store start locally or restore a GitHub backup.
func InitialSetupMenuModel(options *types.Options) MenuModel {
	return MenuModel{
		title: "Password Manager - First-time Setup",
		choices: []MenuItem{
			{Title: "Create a new password store", Description: "Set up a master password and phrase", Action: "setup"},
			{Title: "📥 Restore from GitHub", Description: "Restore a backup and use its original master password", Action: "github_restore"},
			{Title: "🔧 Configure GitHub", Description: "Set or change the backup repository", Action: "github_configure"},
			{Title: "Quit", Description: "Exit the password manager", Action: "quit"},
		},
		options: options,
	}
}
