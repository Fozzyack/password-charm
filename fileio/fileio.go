// Package fileio manages file system operations for the password manager.
// It handles password store initialization, file reading/writing, and directory management
// with appropriate security permissions.
package fileio

import (
	"fmt"
	"log"
	"os"
)

// PasswordFolder represents the password store directory and its current state.
// It tracks the store location, directory contents, initialization status, and master password.
type PasswordFolder struct {
	FolderLocation string        // Absolute path to ~/.password-manager-store/
	Dirs           []os.DirEntry // Contents of the password store directory
	InitCheck      bool          // Whether the store has been properly initialized
	Password       string        // The master password (stored in memory only)
}

// InitPasswordFolder creates or accesses the password store directory and initializes
// the PasswordFolder struct. It creates ~/.password-manager-store/ with secure permissions
// and sets up the .checker subdirectory for validation files.
//
// Returns a fully initialized PasswordFolder or terminates the program on fatal errors.
func InitPasswordFolder() *PasswordFolder {
	passwordFolder := &PasswordFolder{
		InitCheck: true,
	}
	err := getDir(passwordFolder)
	if err != nil {
		log.Fatal("Could not Open or Create Password Store File: ", err)
	}
	return passwordFolder
}

// FileExists checks if a file exists at the given path.
// Returns true if the file exists, false if it doesn't exist or on other errors.
func FileExists(fileName string) bool {
	_, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func getDir(passwordFolder *PasswordFolder) (error) {
	passwordEncFolder := fmt.Sprintf("%s/.password-manager-store", os.Getenv("HOME"))
	dirs, err := os.ReadDir(passwordEncFolder)
	if os.IsNotExist(err) {
		log.Println(".password-manager-store not found\nCreating new encrypted passwords folder")
		os.Mkdir(passwordEncFolder, 0750)
		dirs, err = os.ReadDir(passwordEncFolder)
	} 
	if err != nil {
		log.Fatal("There was an error opening the file\n", err)
		return err
	}

	_, err = os.ReadDir(fmt.Sprintf("%s/.checker", passwordEncFolder)) 
	if os.IsNotExist(err) {
		log.Println("Initialising Checker")
		os.Mkdir(fmt.Sprintf("%s/.checker", passwordEncFolder), 0750)
	} else if err != nil {
		log.Fatal(err)
	}

	if !FileExists(fmt.Sprintf("%s/.checker/init.gpg", passwordEncFolder)) {
		passwordFolder.InitCheck = false
	}


	passwordFolder.Dirs = dirs
	passwordFolder.FolderLocation = passwordEncFolder
	return  nil
}


func (pf *PasswordFolder) WriteToFile (fileName string, input []byte) error {
	err := os.WriteFile(fmt.Sprintf("%s/%s.gpg", pf.FolderLocation, fileName), input, 0666)
	if err != nil {
		log.Printf("ERROR: Writing to file: %s", err)
		return err
	}
	return nil
}

func (pf *PasswordFolder) ReadFromFile (fileName string) ([]byte, error) {
	data, err := os.ReadFile(fmt.Sprintf("%s/%s.gpg", pf.FolderLocation, fileName))
	if err != nil {
		return nil, err
	}
	return data, nil
}

// RefreshDirectoryListing updates the Dirs field with the current contents of the password store directory.
// This is useful when new files have been added since initialization and you need the list to reflect current state.
func (pf *PasswordFolder) RefreshDirectoryListing() error {
	dirs, err := os.ReadDir(pf.FolderLocation)
	if err != nil {
		return fmt.Errorf("failed to read directory %s: %v", pf.FolderLocation, err)
	}
	pf.Dirs = dirs
	return nil
}

// DeleteFile removes a password file from the password store directory.
// The filename should not include the .gpg extension as it will be added automatically.
// Returns an error if the file doesn't exist or if deletion fails.
func (pf *PasswordFolder) DeleteFile(fileName string) error {
	filePath := fmt.Sprintf("%s/%s.gpg", pf.FolderLocation, fileName)
	
	// Check if file exists before attempting deletion
	if !FileExists(filePath) {
		return fmt.Errorf("password file '%s.gpg' does not exist", fileName)
	}
	
	// Attempt to delete the file
	err := os.Remove(filePath)
	if err != nil {
		return fmt.Errorf("failed to delete password file '%s.gpg': %v", fileName, err)
	}
	
	return nil
}






