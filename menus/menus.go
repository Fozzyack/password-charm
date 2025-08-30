// Package menus provides login functionality and user authentication for the password manager.
// It handles first-time setup with validation, master password verification, and login flow management.
package menus

import (
	"fmt"
	"strings"
	"time"
	"github.com/Fozzyack/password-manager/encryption"
	"github.com/Fozzyack/password-manager/fileio"
	"github.com/Fozzyack/password-manager/types"
	"github.com/Fozzyack/password-manager/ui/textinput"
	"github.com/Fozzyack/password-manager/ui/menu"
	"github.com/Fozzyack/password-manager/ui/form"
	"github.com/Fozzyack/password-manager/ui/list"
	"github.com/Fozzyack/password-manager/ui/detail"
	"github.com/Fozzyack/password-manager/ui/confirm"
	"github.com/Fozzyack/password-manager/ui/change"
	"github.com/Fozzyack/password-manager/utils"
	tea "github.com/charmbracelet/bubbletea"
)





type Menu struct {
	passwordFolder      *fileio.PasswordFolder
	encryptionFunctions *encryption.EncryptionFunctions
	Options             *types.Options
}

func InitMenus(pf *fileio.PasswordFolder, ef *encryption.EncryptionFunctions, options *types.Options) *Menu {
	return &Menu{
		passwordFolder:      pf,
		encryptionFunctions: ef,
		Options:             options,
	}
}

// validatePassword checks if password meets minimum length requirement
func validatePassword(password string) (bool, string) {
	const minPasswordLength = 8
	if len(password) < minPasswordLength {
		return false, "Master password must be at least 8 characters long"
	}
	return true, ""
}

// validatePhrase checks if phrase meets minimum length requirement
func validatePhrase(phrase string) (bool, string) {
	const minPhraseLength = 12
	if len(phrase) < minPhraseLength {
		return false, "Phrase must be at least 12 characters long"
	}
	return true, ""
}

func (menu *Menu) Login() (bool, error) {
	// Clear any previous error message before showing login
	menu.Options.ErrorMessage = ""
	
	var err error
	p := tea.NewProgram(textinput.InitialModelWithMasking("Welcome, please type in your Master password", "Password", &menu.passwordFolder.Password, menu.Options, false))

	if !menu.passwordFolder.InitCheck {
		// Validate master password (visible during setup)
		for {
			_, err = p.Run()
			if err != nil {
				return false, err
			}
			if menu.Options.Quit {
				return false, nil
			}
			
			valid, errorMsg := validatePassword(menu.passwordFolder.Password)
			if valid {
				break
			}
			
			// Show validation error and prompt again
			menu.Options.ErrorMessage = errorMsg
			menu.passwordFolder.Password = "" // Clear invalid password
			p = tea.NewProgram(textinput.InitialModelWithMasking("Welcome, please type in your Master password", "Password", &menu.passwordFolder.Password, menu.Options, false))
		}
		
		// Validate phrase
		phrase := ""
		for {
			menu.Options.ErrorMessage = "" // Clear previous errors
			p = tea.NewProgram(textinput.InitialModel("Type in a Random phrase", "the quick brown fox...", &phrase, menu.Options))
			_, err = p.Run()
			if err != nil {
				return false, err
			}
			if menu.Options.Quit {
				return false, nil
			}
			
			valid, errorMsg := validatePhrase(phrase)
			if valid {
				break
			}
			
			// Show validation error
			menu.Options.ErrorMessage = errorMsg
			phrase = "" // Clear invalid phrase
		}
		menu.passwordFolder.InitCheck = false
		data := encryption.Data{
			Password:  phrase,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		err = menu.encryptionFunctions.EncryptPasswordAndWriteToFile(".checker/init", data)
		if err != nil {
			return false, err
		}
		menu.passwordFolder.InitCheck = true
		menu.passwordFolder.Password = ""
	}
	
	p = tea.NewProgram(textinput.InitialModel("Hello Again! Please enter your Password", "Password", &menu.passwordFolder.Password, menu.Options))
	_, err = p.Run(); if err != nil {
		return false, err
	}
	data, err := menu.encryptionFunctions.DecryptPasswordFromFile(".checker/init")
	if err != nil {
		return false, nil
	}
	menu.passwordFolder.Password = data.Password
	return true, nil
}

// ShowMainMenu displays the main menu and handles user selection.
// Returns the selected action string and any error that occurred.
func (m *Menu) ShowMainMenu() (string, error) {
	// Clear any previous error messages
	m.Options.ErrorMessage = ""
	
	// Clear screen before showing menu
	fmt.Print("\033[2J\033[H")
	
	// Create and run the main menu
	mainMenu := menu.InitialMenuModel(m.Options)
	p := tea.NewProgram(mainMenu)
	
	finalModel, err := p.Run()
	if err != nil {
		return "", err
	}
	
	// Check if user quit
	if m.Options.Quit {
		return "quit", nil
	}
	
	// Get the selected action
	menuModel := finalModel.(menu.MenuModel)
	return menuModel.GetSelectedAction(), nil
}

// AddNewPassword displays the add password form and handles password creation.
// Returns true if a password was successfully added, false if cancelled or failed.
func (m *Menu) AddNewPassword() (bool, error) {
	// Clear any previous error messages
	m.Options.ErrorMessage = ""
	
	// Create and run the password form
	passwordForm := form.NewPasswordForm(m.Options)
	p := tea.NewProgram(passwordForm)
	
	finalModel, err := p.Run()
	if err != nil {
		return false, fmt.Errorf("error running form: %v", err)
	}
	
	formModel := finalModel.(form.FormModel)
	
	// Check if form was cancelled
	if formModel.IsCancelled() {
		return false, nil // Not an error, just cancelled
	}
	
	// Check if form was submitted successfully
	if !formModel.IsSubmitted() {
		return false, nil // Form not completed
	}
	
	// Get form data
	formData := formModel.GetFormData()
	
	// Sanitize inputs
	siteName := utils.SanitizeInput(formData["site_service_name"])
	username := utils.SanitizeInput(formData["username"])
	email := utils.SanitizeInput(formData["email"])
	url := utils.SanitizeInput(formData["url"])
	password := formData["password"] // Don't sanitize password to preserve special chars
	
	// Validate required fields
	if siteName == "" || password == "" {
		m.Options.ErrorMessage = "Site name and password are required"
		return false, nil
	}
	
	// Create password entry
	now := time.Now()
	passwordEntry := encryption.Data{
		Password:  password,
		Username:  username,
		Email:     email,
		URL:       url,
		CreatedAt: now,
		UpdatedAt: now,
	}
	
	// Generate filename
	filename := utils.GenerateFilename(siteName)
	
	// Encrypt and save
	err = m.encryptionFunctions.EncryptPasswordAndWriteToFile(filename, passwordEntry)
	if err != nil {
		return false, fmt.Errorf("failed to save password: %v", err)
	}
	
	// Show success message
	fmt.Print("\033[2J\033[H") // Clear screen
	fmt.Printf("✅ Password saved successfully!\n\n")
	fmt.Printf("Site: %s\n", siteName)
	if username != "" {
		fmt.Printf("Username: %s\n", username)
	}
	if email != "" {
		fmt.Printf("Email: %s\n", email)
	}
	if url != "" {
		fmt.Printf("URL: %s\n", url)
	}
	fmt.Printf("File: %s.gpg\n\n", filename)
	
	fmt.Println("Press Enter to continue...")
	fmt.Scanln()
	
	return true, nil
}

// ListAllPasswords displays a list of all stored passwords and handles password viewing.
// Returns true if a password was viewed, false if cancelled or if no passwords exist.
func (m *Menu) ListAllPasswords() (bool, error) {
	// Clear any previous error messages
	m.Options.ErrorMessage = ""
	
	// Get all password entries
	entries, err := m.getAllPasswordEntries()
	if err != nil {
		return false, fmt.Errorf("failed to load password entries: %v", err)
	}
	
	// If no passwords exist, show empty state and return
	if len(entries) == 0 {
		passwordList := list.NewPasswordList(entries, m.Options)
		p := tea.NewProgram(passwordList)
		_, err := p.Run()
		return false, err
	}
	
	// Show the password list
	for {
		passwordList := list.NewPasswordList(entries, m.Options)
		p := tea.NewProgram(passwordList)
		
		finalModel, err := p.Run()
		if err != nil {
			return false, fmt.Errorf("error running password list: %v", err)
		}
		
		listModel := finalModel.(list.ListModel)
		
		// Check if user selected an entry
		if !listModel.IsSelected() {
			// User cancelled or quit
			return false, nil
		}
		
		// Get the selected entry and show details
		selectedEntry := listModel.GetSelectedEntry()
		
		// Decrypt the full password entry
		passwordData, err := m.encryptionFunctions.DecryptPasswordFromFile(selectedEntry.Filename)
		if err != nil {
			fmt.Print("\033[2J\033[H") // Clear screen
			fmt.Printf("❌ Error loading password details: %v\n\n", err)
			fmt.Println("Press Enter to continue...")
			fmt.Scanln()
			continue // Go back to the list
		}
		
		// Show password details
		detailView := detail.NewPasswordDetail(passwordData, selectedEntry.Filename, selectedEntry.SiteName, m.Options)
		detailProgram := tea.NewProgram(detailView)
		
		finalDetailModel, err := detailProgram.Run()
		if err != nil {
			return false, fmt.Errorf("error running password detail view: %v", err)
		}
		
		detailModel := finalDetailModel.(detail.DetailModel)
		
		// Check if deletion was requested
		if detailModel.IsDeletionRequested() {
			// Show confirmation dialog
			confirmDialog := confirm.NewConfirmDialog(selectedEntry.SiteName, selectedEntry.Filename, "delete", m.Options)
			confirmProgram := tea.NewProgram(confirmDialog)
			
			finalConfirmModel, err := confirmProgram.Run()
			if err != nil {
				return false, fmt.Errorf("error running confirmation dialog: %v", err)
			}
			
			confirmModel := finalConfirmModel.(confirm.ConfirmModel)
			
			// Check if deletion was confirmed
			if confirmModel.IsConfirmed() {
				// Delete the file
				err = m.passwordFolder.DeleteFile(selectedEntry.Filename)
				if err != nil {
					fmt.Print("\033[2J\033[H") // Clear screen
					fmt.Printf("❌ Error deleting password: %v\n\n", err)
					fmt.Println("Press Enter to continue...")
					fmt.Scanln()
					continue // Return to list
				}
				
				// Refresh directory listing and entries after deletion
				err = m.passwordFolder.RefreshDirectoryListing()
				if err != nil {
					fmt.Print("\033[2J\033[H") // Clear screen
					fmt.Printf("❌ Error refreshing password list: %v\n\n", err)
					fmt.Println("Press Enter to continue...")
					fmt.Scanln()
					continue // Return to list anyway
				}
				
				// Get updated entries
				entries, err = m.getAllPasswordEntries()
				if err != nil {
					return false, fmt.Errorf("failed to reload password entries after deletion: %v", err)
				}
				
				// Show success message
				fmt.Print("\033[2J\033[H") // Clear screen
				fmt.Printf("✅ Password deleted successfully!\n\n")
				fmt.Printf("Deleted: %s (%s.gpg)\n\n", selectedEntry.SiteName, selectedEntry.Filename)
				fmt.Println("Press Enter to continue...")
				fmt.Scanln()
				
				// Continue to show updated list (entries variable is already updated)
				continue
			}
			// If deletion was cancelled, return to detail view for the same entry
			continue
		}
		
		// After viewing details without deletion, return to the list (continue the loop)
		// User can press Esc from the list to exit completely
	}
}

// getAllPasswordEntries retrieves and decrypts all password entries from the store
func (m *Menu) getAllPasswordEntries() ([]list.PasswordEntry, error) {
	var entries []list.PasswordEntry
	
	// Refresh directory listing to ensure we have the latest files
	err := m.passwordFolder.RefreshDirectoryListing()
	if err != nil {
		return nil, fmt.Errorf("failed to refresh directory listing: %v", err)
	}
	
	// Iterate through all files in the password store
	for _, dirEntry := range m.passwordFolder.Dirs {
		// Skip directories and non-GPG files
		if dirEntry.IsDir() || !strings.HasSuffix(dirEntry.Name(), ".gpg") {
			continue
		}
		
		// Skip the init.gpg file used for validation
		if dirEntry.Name() == "init.gpg" {
			continue
		}
		
		// Skip files in .checker directory
		if strings.Contains(dirEntry.Name(), ".checker") {
			continue
		}
		
		// Get filename without .gpg extension
		filename := strings.TrimSuffix(dirEntry.Name(), ".gpg")
		
		// Try to decrypt the entry to get its details
		passwordData, err := m.encryptionFunctions.DecryptPasswordFromFile(filename)
		if err != nil {
			// If we can't decrypt a file, skip it but don't fail entirely
			// This allows the user to see other passwords even if one is corrupted
			continue
		}
		
		// Parse the site name from filename
		siteName := utils.ParseFilenameToSiteName(filename)
		
		// Create list entry
		entry := list.PasswordEntry{
			Filename:  filename,
			SiteName:  siteName,
			Username:  passwordData.Username,
			Email:     passwordData.Email,
			CreatedAt: passwordData.CreatedAt,
		}
		
		entries = append(entries, entry)
	}
	
	return entries, nil
}

// ChangeMasterPassword handles the master password change workflow.
// Returns true if password was changed successfully, false if cancelled or failed.
func (m *Menu) ChangeMasterPassword() (bool, error) {
	// Clear any previous error messages
	m.Options.ErrorMessage = ""

	// Show the change password form
	changeForm := change.NewChangePasswordForm(m.Options)
	p := tea.NewProgram(changeForm)
	
	finalModel, err := p.Run()
	if err != nil {
		return false, fmt.Errorf("error running change password form: %v", err)
	}
	
	formModel := finalModel.(change.ChangeModel)
	
	// Check if form was cancelled
	if formModel.IsCancelled() {
		return false, nil // Not an error, just cancelled
	}
	
	// Check if form was submitted successfully
	if !formModel.IsSubmitted() {
		return false, nil // Form not completed
	}
	
	// Get form data
	currentPass, newPass, _ := formModel.GetFormData()
	
	// Step 1: Verify current password by trying to decrypt init.gpg
	oldPassword := m.passwordFolder.Password // Store original password
	m.passwordFolder.Password = currentPass  // Temporarily set to verify
	
	validationData, err := m.encryptionFunctions.DecryptPasswordFromFile(".checker/init")
	if err != nil {
		// Restore original password
		m.passwordFolder.Password = oldPassword
		
		fmt.Print("\033[2J\033[H") // Clear screen
		fmt.Printf("❌ Error: Current password is incorrect\n\n")
		fmt.Println("Press Enter to continue...")
		fmt.Scanln()
		return false, nil
	}
	
	// Step 2: Re-encrypt the validation data with new master password
	m.passwordFolder.Password = newPass // Set new password for encryption
	
	// Update the timestamp to reflect the password change
	validationData.UpdatedAt = time.Now()
	
	err = m.encryptionFunctions.EncryptPasswordAndWriteToFile(".checker/init", validationData)
	if err != nil {
		// Restore original password on failure
		m.passwordFolder.Password = oldPassword
		
		fmt.Print("\033[2J\033[H") // Clear screen
		fmt.Printf("❌ Error saving new master password: %v\n\n", err)
		fmt.Println("Press Enter to continue...")
		fmt.Scanln()
		return false, nil
	}
	
	// Step 3: Test that we can decrypt with the new password
	testData, err := m.encryptionFunctions.DecryptPasswordFromFile(".checker/init")
	if err != nil {
		// This shouldn't happen, but if it does, we're in trouble
		fmt.Print("\033[2J\033[H") // Clear screen
		fmt.Printf("❌ Critical Error: Cannot decrypt with new password. Please check your password store manually.\n\n")
		fmt.Println("Press Enter to continue...")
		fmt.Scanln()
		return false, fmt.Errorf("critical error: new password verification failed: %v", err)
	}
	
	// Verify the validation phrase is still correct
	if testData.Password != validationData.Password {
		fmt.Print("\033[2J\033[H") // Clear screen
		fmt.Printf("❌ Critical Error: Validation data corrupted during password change.\n\n")
		fmt.Println("Press Enter to continue...")
		fmt.Scanln()
		return false, fmt.Errorf("validation data integrity check failed")
	}
	
	// Step 4: Success! Show confirmation message
	fmt.Print("\033[2J\033[H") // Clear screen
	fmt.Printf("✅ Master password changed successfully!\n\n")
	fmt.Printf("Your new master password is now active.\n")
	fmt.Printf("You will need to use the new password for future logins.\n\n")
	fmt.Println("Press Enter to continue...")
	fmt.Scanln()
	
	// The new password is already set in m.passwordFolder.Password
	// so the current session continues to work normally
	
	return true, nil
}
