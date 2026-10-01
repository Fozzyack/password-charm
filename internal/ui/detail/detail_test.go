package detail

import (
	"errors"
	"strings"
	"testing"

	"github.com/Fozzyack/password-manager/internal/encryption"
	tea "github.com/charmbracelet/bubbletea"
)

func TestCopyPasswordKeepsPasswordHidden(t *testing.T) {
	model := DetailModel{entry: encryption.Data{Password: "test-secret-123!"}}
	updated, cmd := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'c'}})
	model = updated.(DetailModel)
	if cmd == nil {
		t.Fatal("copy action should schedule a clipboard write")
	}
	if model.IsPasswordVisible() {
		t.Fatal("copy action should not reveal the password")
	}
	if strings.Contains(model.View(), model.entry.Password) {
		t.Fatal("copy feedback should not expose the password")
	}
	if !strings.Contains(model.View(), "Copying password...") {
		t.Fatal("copy action should display pending feedback")
	}
	_, repeatedCmd := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'c'}})
	if repeatedCmd != nil {
		t.Fatal("repeated copy requests should wait for the pending write")
	}
}

func TestClipboardFeedbackAndRetry(t *testing.T) {
	for _, test := range []struct {
		name string
		err  error
		want string
	}{
		{name: "success", want: "Password copied to clipboard"},
		{name: "failure", err: errors.New("clipboard unavailable"), want: "Could not copy password: clipboard unavailable"},
	} {
		t.Run(test.name, func(t *testing.T) {
			model := DetailModel{copying: true, entry: encryption.Data{Password: "test-secret-123!"}}
			updated, _ := model.Update(clipboardResultMsg{err: test.err})
			model = updated.(DetailModel)
			if !strings.Contains(model.View(), test.want) {
				t.Fatalf("expected clipboard feedback %q", test.want)
			}
			if strings.Contains(model.View(), model.entry.Password) {
				t.Fatal("clipboard result should not expose the password")
			}
			_, cmd := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'c'}})
			if cmd == nil {
				t.Fatal("copy should be available again after the result")
			}
		})
	}
}
