// Package main provides the entry point for the GPG-based password manager application.
// This application uses terminal UI for secure password management with local GPG encryption.
package main

import (
	"fmt"

	"github.com/Fozzyack/password-manager/encryption"
	"github.com/Fozzyack/password-manager/fileio"
	"github.com/Fozzyack/password-manager/menus"
	"github.com/Fozzyack/password-manager/types"
)

// main is the application entry point. It initializes the password store,
// handles the login flow with validation, and manages the main application loop.
// The application supports first-time setup with master password and phrase validation,
// as well as subsequent logins with password verification.
func main() {

	passwordFolder := fileio.InitPasswordFolder()
	var err error

	for file := range(passwordFolder.Dirs) {
		fmt.Println(passwordFolder.Dirs[file])
	}
	fmt.Println(passwordFolder.InitCheck)
	options := &types.Options{
		Quit : false,
		LoggedIn: false,
		ErrorMessage: "",
	}
	encrypt := encryption.NewEncryption(passwordFolder)
	menu := menus.InitMenus(passwordFolder, encrypt, options)
	for !options.LoggedIn && !options.Quit{
		options.LoggedIn, err = menu.Login()
		if err != nil {
			panic(err)
		} else if options.Quit {
			fmt.Print("\033[2J\033[H") // Clear screen
			fmt.Println("Escape Sequence Detected :: Exiting")
		} else if !options.LoggedIn {
			// Clear screen and show error message
			fmt.Print("\033[2J\033[H")
			fmt.Printf("\033[91m\033[1mIncorrect Password - Try again\033[0m\n\n")
			fmt.Println("Press Enter to continue...")
			fmt.Scanln() // Wait for user to press Enter
		} else {
			// Clear error message on successful login
			options.ErrorMessage = ""
		}
	}

	// Main menu loop after successful login
	for options.LoggedIn && !options.Quit {
		action, err := menu.ShowMainMenu()
		if err != nil {
			fmt.Printf("Error displaying menu: %v\n", err)
			break
		}

		// Handle the selected action
		handleMenuAction(action, menu)
		
		// Check if user wants to quit
		if action == "quit" || options.Quit {
			fmt.Print("\033[2J\033[H") // Clear screen
			fmt.Println("Goodbye! 👋")
			break
		}
	}
}

// handleMenuAction processes the selected menu action and calls appropriate functions
func handleMenuAction(action string, menu *menus.Menu) {
	fmt.Print("\033[2J\033[H") // Clear screen

	switch action {
	case "list":
		_, err := menu.ListAllPasswords()
		if err != nil {
			fmt.Print("\033[2J\033[H") // Clear screen
			fmt.Printf("❌ Error listing passwords: %v\n\n", err)
			waitForEnter()
		}

	case "add":
		_, err := menu.AddNewPassword()
		if err != nil {
			fmt.Print("\033[2J\033[H") // Clear screen
			fmt.Printf("❌ Error adding password: %v\n\n", err)
			waitForEnter()
		}
		// Success message is handled within AddNewPassword

	case "change_master":
		fmt.Println("🔄 Changing master password...")
		fmt.Println("This feature is coming soon!")
		waitForEnter()

	case "export":
		fmt.Println("📤 Exporting passwords...")
		fmt.Println("This feature is coming soon!")
		waitForEnter()

	case "quit":
		// Handled in main loop
		return

	default:
		fmt.Printf("Unknown action: %s\n", action)
		waitForEnter()
	}
}

// waitForEnter pauses execution until the user presses Enter
func waitForEnter() {
	fmt.Println("\nPress Enter to continue...")
	fmt.Scanln()
}
