# Password Charm

Password Charm is a personal project developed to simplify local password management.

A terminal-based password manager built with [Bubble Tea](https://github.com/charmbracelet/bubbletea), Bubbles, and Lip Gloss. Password entries are stored locally, encrypted using password-based OpenPGP through ProtonMail's [GopenPGP](https://github.com/ProtonMail/gopenpgp) library. Optional GitHub integration provides encrypted backups and restoration on a new machine.

## Features

- **Local encrypted storage** - passwords and entry metadata are encrypted before being written to disk
- **Keyboard-driven interface** - navigate menus, forms, and password entries from the terminal
- **Add new passwords** - save a site/service name and password, with optional username, email, and URL
- **View your passwords** - browse entries, reveal passwords, and see a basic password-strength rating and creation timestamp
- **Copy passwords** - copy a selected password to the system clipboard without revealing it
- **Delete entries** - confirmation is required before an entry is removed
- **Change master password** - verify your current password and set a new login password
- **Master password authentication** - access the password store with a single login password
- **GitHub backup and restore** - save an encrypted snapshot to a GitHub repository and restore it into an empty password store

## Getting Started

### Requirements

- **Go 1.24.6 or later** to build from source
- **Git** to clone the repository and use GitHub backup or restore
- A terminal that supports ANSI colors and keyboard input

Encryption runs inside the application; a separate GPG installation or key pair is not required.

On Linux, clipboard copying requires `wl-clipboard` for Wayland or `xclip` / `xsel` for X11. The application reports a copy error if a clipboard utility is unavailable.

### Build and Run

1. **Clone and build**:
   ```bash
   git clone https://github.com/Fozzyack/password-charm.git
   cd password-charm
   go build -o password-manager ./cmd/password-manager
   ```
2. **Run the application**:
   ```bash
   ./password-manager
   ```

You can also run from the repository with `go run ./cmd/password-manager`. Go downloads the required dependencies on the first build or run.

## Usage

### Initial Setup

1. Choose **Create a new password store** from the first-time setup menu. To use an existing backup instead, follow the restore instructions below.
2. Create a master password of at least 8 characters.
3. Create a random validation phrase of at least 12 characters. The application encrypts this phrase to verify your master password on subsequent logins.
4. Enter your master password again to log in.
5. Choose **Add New Password** to save your first entry.

### Daily Use

1. Enter your master password to log in.
2. Choose **Add New Password** to enter a site/service name and password. Username, email, and URL are optional. Press Enter on the final field to save.
3. Choose **List All Passwords**, then select an entry to see its details. Passwords are hidden by default; press `v` or Space to reveal or hide them.
4. Press `c` in the detail view to copy the password to the clipboard. The screen displays whether the copy succeeded.
5. Press `d` in the detail view to request deletion, then confirm or cancel.

## GitHub Backup and Restore

### Requirements

- An existing GitHub repository, preferably a dedicated private repository
- Git authentication configured through SSH keys or an HTTPS credential helper
- An internet connection for backup and restore operations

The application uses your existing Git credentials and does not store GitHub tokens. Authentication must be configured before starting a backup or restore; Git's interactive credential prompts are disabled.

### Configure a Repository

Choose **Configure GitHub** from the main menu or first-time setup menu and enter a GitHub repository URL, such as:

```text
git@github.com:owner/repository.git
https://github.com/owner/repository.git
ssh://git@github.com/owner/repository.git
```

URLs must not contain embedded credentials. The application saves the repository URL in `password-manager/github.json` within your operating system's user configuration directory. On Linux, this is normally `~/.config/password-manager/github.json`, or under `$XDG_CONFIG_HOME` when set.

### Save a Backup

1. Log in to your password store.
2. Choose **Save to GitHub** from the main menu. If no repository is configured, the application prompts for one.
3. Wait for the operation to complete and review the result.

Each backup updates the repository's `passwords/` directory with a complete snapshot of the encrypted `.gpg` entries and the encrypted login checker. Deleted entries are removed from the latest snapshot. Changes are committed and pushed to the repository's default branch; an empty repository uses `main`. If the snapshot has not changed, no new commit is created.

Backups run only when **Save to GitHub** is selected; local changes are not uploaded automatically.

### Restore a Backup

1. Start the application on a machine with an empty password store.
2. Choose **Restore from GitHub** from the first-time setup menu.
3. Enter the repository URL if prompted, or use **Configure GitHub** to change the saved repository first.
4. After restoration completes, log in with the master password associated with the backup.

Restore installs the latest encrypted snapshot without overwriting existing entries. It requires an empty store and is offered before a new store's master password and validation phrase are created.

## Keyboard Shortcuts

| Screen | Keys | Action |
| --- | --- | --- |
| Menus and password list | ↑/↓ or `k`/`j` | Move between items |
| Menus and password list | Enter or Space | Select an item |
| Password list | Home / End | Jump to the first / last entry |
| Forms | Tab / Shift+Tab or ↓/↑ | Move between fields |
| Forms | Enter | Move to the next field; submit on the final field |
| Password details | `v` or Space | Reveal or hide the password |
| Password details | `c` | Copy the password to the clipboard |
| Password details | `d` | Request deletion |
| Password details | Esc, `q`, Backspace, or Enter | Return to the list |
| Delete confirmation | `y` / `n` | Confirm / cancel deletion |
| Forms and confirmation | Esc or Ctrl+C | Cancel the current action |
| Password list | Esc, `q`, or Ctrl+C | Return to the main menu |
| Main menu and login | Esc or Ctrl+C | Quit the application |

## Planned Features

- **Password export** - export entries in CSV or JSON format

Password export is not yet implemented.

## Encryption and Storage

- Password data is stored locally; local password management does not require an internet connection
- Uses password-based OpenPGP encryption with GopenPGP's RFC 9580 profile
- Passwords, usernames, email addresses, URLs, and timestamps are encrypted together in ASCII-armored `.gpg` files
- Site/service names are used to generate filenames and remain visible on disk
- GitHub backups contain encrypted entries and the encrypted login checker; entry filenames also remain visible in the backup repository
- Your password store is created automatically at `~/.password-manager-store/`
