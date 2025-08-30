// Package encryption provides GPG-based encryption and decryption functionality for the password manager.
// It uses ProtonMail's gopenpgp library with RFC9580 OpenPGP profile for secure password storage.
package encryption

import (
	"encoding/json"
	"time"

	"github.com/ProtonMail/gopenpgp/v3/crypto"
	"github.com/ProtonMail/gopenpgp/v3/profile"

	"github.com/Fozzyack/password-manager/fileio"
)

// Data represents a password entry with associated metadata.
// All fields are JSON-serialized before encryption for secure storage.
type Data struct {
	Password  string    `json:"password"`  // The actual password or secret data
	Username  string    `json:"username"`  // Associated username (optional)
	Email     string    `json:"email"`     // Associated email address (optional)
	URL       string    `json:"url"`       // Associated website URL (optional)
	CreatedAt time.Time `json:"created_at"` // Timestamp when entry was created
	UpdatedAt time.Time `json:"updated_at"` // Timestamp when entry was last modified
}

// EncryptionFunctions provides methods for encrypting and decrypting password data
// using the master password from the password folder.
type EncryptionFunctions struct {
	passwordFolder  *fileio.PasswordFolder // Reference to the password store
	EnteredPassword string                  // Currently unused, may be removed
}

// NewEncryption creates a new EncryptionFunctions instance with the given password folder.
// The password folder must be initialized and contain the master password for encryption operations.
func NewEncryption(passwordFolder *fileio.PasswordFolder) *EncryptionFunctions {
	return &EncryptionFunctions{
		passwordFolder: passwordFolder,
	}
}

// EncryptPasswordAndWriteToFile encrypts the given Data struct and writes it to a file.
// The data is first JSON-serialized, then encrypted using the master password with GPG,
// and finally written as an armored .gpg file in the password store.
//
// Parameters:
//   - fileName: Name of the file (without .gpg extension, added automatically)
//   - data: Data struct containing password and metadata to encrypt
//
// Returns an error if JSON marshaling, encryption, or file writing fails.
func (ef *EncryptionFunctions) EncryptPasswordAndWriteToFile(fileName string, data Data) error {
	// Convert the Data struct to JSON for storage
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// Use the master password for encryption
	password := []byte(ef.passwordFolder.Password)
	pgp := crypto.PGPWithProfile(profile.RFC9580())

	// Create encryption handler with password-based encryption
	encHandle, err := pgp.Encryption().Password(password).New()
	if err != nil {
		return err
	}

	// Encrypt the JSON data
	pgpMessage, err := encHandle.Encrypt(jsonData)
	if err != nil {
		return err
	}

	// Convert to ASCII-armored format for storage
	armored, err := pgpMessage.ArmorBytes()
	if err != nil {
		return err
	}

	// Write the encrypted data to file (adds .gpg extension automatically)
	err = ef.passwordFolder.WriteToFile(fileName, armored)
	if err != nil {
		return err
	}

	return nil
}

func (ef *EncryptionFunctions) DecryptPasswordFromFile (fileName string) (Data, error) {

	data := Data{}
	fileData, err := ef.passwordFolder.ReadFromFile(fileName)
	if err != nil {
		return data, err 
	}
	password := []byte(ef.passwordFolder.Password)
	pgp := crypto.PGPWithProfile(profile.RFC9580())

	decHandler, err := pgp.Decryption().Password(password).New()
	if err != nil {
		return data, err
	}
	decrypted, err := decHandler.Decrypt(fileData, crypto.Armor)
	if err != nil {
		return data, err
	}
	json.Unmarshal(decrypted.Bytes(), &data)
	return data, nil
}
