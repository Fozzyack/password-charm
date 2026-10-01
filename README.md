# Password Charm

> **Personal Note**: This is a hobby project born out of frustration - I can never remember my passwords! 🤦‍♂️

A terminal-based password manager built with [Bubble Tea](https://github.com/charmbracelet/bubbletea), Bubbles, and Lip Gloss. Password entries stay on your machine, encrypted using password-based OpenPGP through ProtonMail's [GopenPGP](https://github.com/ProtonMail/gopenpgp) library.

## ✨ What it does

- **Local encrypted storage** - passwords and entry metadata are encrypted before being written to disk
- **Easy to use** - simple keyboard navigation through menus
- **Add new passwords** - save a site/service name and password, with optional username, email, and URL
- **View your passwords** - browse entries, reveal passwords, and see a basic password-strength rating and creation timestamp
- **Delete old passwords** - with confirmation to prevent accidents
- **Change master password** - verify your current password and set a new login password
- **Master password protection** - one password to access everything

## 🚀 Quick Start

### Requirements

- **Go 1.24.6 or later** to build from source
- **Git** to clone the repository
- A terminal that supports ANSI colors and keyboard input

Encryption runs inside the application; a separate GPG installation or key pair is not required.

### Build and run

1. **Clone and build**:
   ```bash
   git clone https://github.com/Fozzyack/password-charm.git
   cd password-charm
   go build -o password-manager .
   ```
2. **Run it**:
   ```bash
   ./password-manager
   ```

You can also run from the repository with `go run .`. Go downloads the required dependencies on the first build or run.

## 📖 How to Use

### First Time
1. Create a master password (8+ characters)
2. Create a random validation phrase (12+ characters). The app encrypts this phrase to check your master password on future logins.
3. Enter your master password again to log in.
4. Choose **Add New Password** to save your first entry.

### Daily Use
1. Enter your master password
2. Choose **Add New Password** to enter a site/service name and password. Username, email, and URL are optional. Press Enter on the final field to save.
3. Choose **List All Passwords**, then select an entry to see its details. Passwords are hidden by default; press `v` or Space to reveal or hide them.
4. Press `d` in the detail view to request deletion, then confirm or cancel.

## ⌨️ Keyboard Shortcuts

| Screen | Keys | Action |
| --- | --- | --- |
| Main menu and password list | ↑/↓ or `k`/`j` | Move between items |
| Main menu and password list | Enter or Space | Select an item |
| Password list | Home / End | Jump to the first / last entry |
| Forms | Tab / Shift+Tab or ↓/↑ | Move between fields |
| Forms | Enter | Move to the next field; submit on the final field |
| Password details | `v` or Space | Reveal or hide the password |
| Password details | `d` | Request deletion |
| Password details | Esc, `q`, Backspace, or Enter | Return to the list |
| Delete confirmation | `y` / `n` | Confirm / cancel deletion |
| Forms and confirmation | Esc or Ctrl+C | Cancel the current action |
| Password list | Esc, `q`, or Ctrl+C | Return to the main menu |
| Main menu and login | Esc or Ctrl+C | Quit the application |


## 📝 To Be Added

- **📤 Export Passwords**: Export to CSV/JSON formats
- **🔗 Maybe: GitHub Backup**: Encrypted backup to private repos

The **Export Passwords** menu option currently displays a coming-soon message.

## 🔒 Security

- Everything stays on your computer (no internet required)
- Uses password-based OpenPGP encryption with GopenPGP's RFC 9580 profile
- Passwords, usernames, email addresses, URLs, and timestamps are encrypted together in ASCII-armored `.gpg` files
- Site/service names are used to generate filenames and remain visible on disk
- Your password store is created automatically at `~/.password-manager-store/`
