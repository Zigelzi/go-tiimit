package auth

import (
	"errors"
	"strings"
	"unicode/utf8"
)

// The limits of the username and password rules for a new account. Lengths
// count letters (runes), not bytes.
const (
	usernameMinLength    = 5
	usernameMaxLength    = 20
	newPasswordMinLength = 10
	newPasswordMaxLength = 64
	// passwordMaxBytes is bcrypt's limit: it ignores every byte after the 72nd.
	passwordMaxBytes = 72
)

// The ways a new username can break the username rule.
var (
	ErrUsernameEmpty   = errors.New("username is empty")
	ErrUsernameLength  = errors.New("username must have 5 to 20 letters")
	ErrUsernameLetters = errors.New("username can have only the letters a to z, å, ä and ö")
)

// The ways a new password can break the password rule.
var (
	ErrPasswordShort          = errors.New("password has fewer than 10 characters")
	ErrPasswordLong           = errors.New("password has more than 64 characters")
	ErrPasswordBytes          = errors.New("password has more than 72 bytes")
	ErrPasswordSameAsUsername = errors.New("password is the same as the username")
)

// NormalizeUsername returns s without surrounding whitespace and in lowercase.
// Both creating an account and logging in use it, so " AinVir " and "ainvir"
// are the same user.
func NormalizeUsername(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

// ValidateUsername returns nil when a normalized username follows the rule for
// a new account: 5 to 20 letters, each one a to z, å, ä or ö. Otherwise it
// returns one of the ErrUsername errors. Login does not check this rule, so
// older usernames keep working.
func ValidateUsername(normalized string) error {
	if normalized == "" {
		return ErrUsernameEmpty
	}
	length := utf8.RuneCountInString(normalized)
	if length < usernameMinLength || length > usernameMaxLength {
		return ErrUsernameLength
	}
	for _, r := range normalized {
		if !isUsernameLetter(r) {
			return ErrUsernameLetters
		}
	}
	return nil
}

// isUsernameLetter reports whether r is a letter a username may contain.
func isUsernameLetter(r rune) bool {
	return (r >= 'a' && r <= 'z') || r == 'å' || r == 'ä' || r == 'ö'
}

// ValidatePassword returns nil when password follows the rule for a new
// account: 10 to 64 characters, at most 72 bytes, and not the same as the
// normalized username. Otherwise it returns one of the ErrPassword errors.
// The password is checked as typed and must be stored as typed: it is never
// trimmed or lowercased.
func ValidatePassword(password, normalizedUsername string) error {
	length := utf8.RuneCountInString(password)
	if length < newPasswordMinLength {
		return ErrPasswordShort
	}
	if length > newPasswordMaxLength {
		return ErrPasswordLong
	}
	if len(password) > passwordMaxBytes {
		return ErrPasswordBytes
	}
	if NormalizeUsername(password) == normalizedUsername {
		return ErrPasswordSameAsUsername
	}
	return nil
}
