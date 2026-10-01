package main

import (
	"fmt"

	"github.com/Fozzyack/password-manager/internal/menus"
)

func handleGitHubAction(action string, menu *menus.Menu) {
	var saved bool
	var err error
	if action == "github_configure" {
		saved, err = menu.ConfigureGitHub()
		if err != nil {
			fmt.Printf("Error configuring GitHub: %v\n", err)
		} else if saved {
			fmt.Println("GitHub repository configured.")
		}
	} else {
		saved, err = menu.BackupToGitHub()
		if err != nil {
			fmt.Printf("Error saving to GitHub: %v\n", err)
		} else if saved {
			fmt.Println("Encrypted passwords saved to GitHub.")
		}
	}
	if !menu.Options.Quit {
		waitForEnter()
	}
}
