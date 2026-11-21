package auth

import (
	"errors"
	"strings"
	"testing"
)

// The username and password rules for a new account, from
// spec/milestones/02-user-registration.md (Decided, 02/R5, 02/Q4, 02/R10).

func TestNormalizeUsername(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{"ainvir", "ainvir"},
		{" AinVir ", "ainvir"},
		{"\tAINVIR\n", "ainvir"},
		{"  ÅSAÄÖ  ", "åsaäö"},
		{"", ""},
		{"   ", ""},
	}
	for _, tc := range cases {
		if got := NormalizeUsername(tc.in); got != tc.want {
			t.Errorf("NormalizeUsername(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

// ValidateUsername gets a username that is already normalised.
func TestValidateUsername(t *testing.T) {
	cases := []struct {
		name     string
		username string
		want     error
	}{
		{"five letters", "ainvi", nil},
		{"twenty letters", strings.Repeat("a", 20), nil},
		{"nordic letters", "åsaäö", nil},
		{"empty", "", ErrUsernameEmpty},
		{"four letters", "aino", ErrUsernameLength},
		{"twenty-one letters", strings.Repeat("a", 21), ErrUsernameLength},
		{"four nordic letters count as four", "åäöå", ErrUsernameLength},
		{"dash", "ain-vir", ErrUsernameLetters},
		{"digit", "ainvir2", ErrUsernameLetters},
		{"space inside", "ain vir", ErrUsernameLetters},
		{"letter outside a to z and å ä ö", "renée", ErrUsernameLetters},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateUsername(tc.username)
			if tc.want == nil && err != nil {
				t.Fatalf("ValidateUsername(%q) = %v, want no error", tc.username, err)
			}
			if tc.want != nil && !errors.Is(err, tc.want) {
				t.Fatalf("ValidateUsername(%q) = %v, want %v", tc.username, err, tc.want)
			}
		})
	}
}

// ValidatePassword gets the password as typed and the normalised username.
func TestValidatePassword(t *testing.T) {
	cases := []struct {
		name     string
		password string
		username string
		want     error
	}{
		{"ten characters", "kymmenen10", "ainvir", nil},
		{"spaces and capitals kept", "Kesä Treeni 2026", "ainvir", nil},
		{"sixty-four characters", strings.Repeat("a", 64), "ainvir", nil},
		{"ten nordic letters count as ten, not twenty bytes", strings.Repeat("ä", 10), "ainvir", nil},
		{"exactly 72 bytes", strings.Repeat("ä", 36), "ainvir", nil},
		{"nine characters", "yhdeksän9", "ainvir", ErrPasswordShort},
		{"nine nordic letters", strings.Repeat("ö", 9), "ainvir", ErrPasswordShort},
		{"empty", "", "ainvir", ErrPasswordShort},
		{"sixty-five characters", strings.Repeat("a", 65), "ainvir", ErrPasswordLong},
		{"sixty-four characters but over 72 bytes", strings.Repeat("ä", 37) + strings.Repeat("a", 27), "ainvir", ErrPasswordBytes},
		{"same as the username", "ainovirtanen", "ainovirtanen", ErrPasswordSameAsUsername},
		{"same as the username in other capitals", "AinoVirtanen", "ainovirtanen", ErrPasswordSameAsUsername},
		{"same as the username with spaces around", " AinoVirtanen ", "ainovirtanen", ErrPasswordSameAsUsername},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidatePassword(tc.password, tc.username)
			if tc.want == nil && err != nil {
				t.Fatalf("ValidatePassword(%q, %q) = %v, want no error", tc.password, tc.username, err)
			}
			if tc.want != nil && !errors.Is(err, tc.want) {
				t.Fatalf("ValidatePassword(%q, %q) = %v, want %v", tc.password, tc.username, err, tc.want)
			}
		})
	}
}
