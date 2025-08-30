# Password Manager

> **Personal Note**: This is a hobby project born out of frustration - I can never remember my passwords! 🤦‍♂️

A secure, terminal-based password manager with a beautiful interface. All your passwords stay on your machine, encrypted with GPG.

## ✨ What it does

- **Store passwords securely** - everything encrypted and stored locally
- **Easy to use** - simple keyboard navigation through menus
- **Add new passwords** - fill out forms for websites/services  
- **View your passwords** - browse and reveal passwords when needed
- **Delete old passwords** - with confirmation to prevent accidents
- **Master password protection** - one password to access everything

## 🚀 Quick Start

1. **Install Go** (1.24.6 or later)
2. **Clone and build**:
   ```bash
   git clone <repository-url>
   cd password-manager
   go mod tidy
   go build -o password-manager
   ```
3. **Run it**:
   ```bash
   ./password-manager
   ```

## 📖 How to Use

### First Time
1. Create a master password (8+ characters)
2. Create a validation phrase (12+ characters)
3. You're ready to store passwords!

### Daily Use
1. Enter your master password
2. Use the menu to add, view, or delete passwords
3. Press 'd' in password details to delete (with confirmation)
4. Press 'v' to show/hide passwords when viewing

## ⌨️ Keyboard Shortcuts

- **Arrow keys / j/k**: Navigate menus and lists
- **Enter/Space**: Select items or confirm actions
- **v**: Show/hide passwords when viewing
- **d**: Delete password (asks for confirmation)
- **Esc**: Go back or cancel
- **Ctrl+C**: Quit application


## 📝 To Be Added

- **🔄 Change Master Password**: Update your master password  
- **📤 Export Passwords**: Export to CSV/JSON formats
- **🔗 Maybe: GitHub Backup**: Encrypted backup to private repos

## 🔒 Security

- Everything stays on your computer (no internet required)
- Uses GPG encryption (battle-tested security)
- Your passwords are stored in `~/.password-manager-store/`
