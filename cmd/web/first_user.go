package main

import (
	"context"
	"fmt"

	"github.com/Zigelzi/go-tiimit/internal/auth"
	"github.com/Zigelzi/go-tiimit/internal/db"
)

// The settings that create the first user of an empty database.
const (
	firstUserUsernameSetting = "FIRST_USER_USERNAME"
	firstUserPasswordSetting = "FIRST_USER_PASSWORD"
)

// firstUserOutcome is what ensureFirstUser did when it returned no error.
type firstUserOutcome int

const (
	// firstUserCreated: the database was empty and the first user was created.
	firstUserCreated firstUserOutcome = iota
	// firstUserNotNeeded: the database already has users; the settings were ignored.
	firstUserNotNeeded
	// firstUserNotSet: the database is empty and neither setting is set, so
	// nobody can log in yet.
	firstUserNotSet
)

// ensureFirstUser creates the first user of an empty database from the
// first-user settings (milestone 03). It counts users before it looks at the
// settings, so on a database with users the settings are ignored even when
// they break the rules, and a leftover setting cannot stop the live app
// (03/R1). On an empty database, a half-set setting or one that breaks the
// sign-up rules returns an error naming the setting; the error never contains
// the password.
//
// The username is normalized as at sign-up and the password is used exactly
// as given. The count and the create run in one transaction, and SQLite lets
// only one of two racing transactions write, so two starts at once create at
// most one first user (03/R10).
func (cfg *webConfig) ensureFirstUser(ctx context.Context, username, password string) (firstUserOutcome, error) {
	tx, err := cfg.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("starting the transaction: %w", err)
	}
	defer tx.Rollback()
	queries := cfg.queries.WithTx(tx)

	count, err := queries.CountUsers(ctx)
	if err != nil {
		return 0, fmt.Errorf("counting users: %w", err)
	}
	if count > 0 {
		return firstUserNotNeeded, nil
	}

	username = auth.NormalizeUsername(username)
	switch {
	case username == "" && password == "":
		return firstUserNotSet, nil
	case username == "":
		return 0, fmt.Errorf("%s is missing: set it together with %s", firstUserUsernameSetting, firstUserPasswordSetting)
	case password == "":
		return 0, fmt.Errorf("%s is missing: set it together with %s", firstUserPasswordSetting, firstUserUsernameSetting)
	}
	if err := auth.ValidateUsername(username); err != nil {
		return 0, fmt.Errorf("%s breaks the sign-up rules: %w", firstUserUsernameSetting, err)
	}
	if err := auth.ValidatePassword(password, username); err != nil {
		return 0, fmt.Errorf("%s breaks the sign-up rules: %w", firstUserPasswordSetting, err)
	}

	hashed, err := auth.HashPassword(password)
	if err != nil {
		return 0, fmt.Errorf("hashing the first user's password: %w", err)
	}
	_, err = queries.CreateUser(ctx, db.CreateUserParams{
		Username:       username,
		HashedPassword: hashed,
	})
	if err != nil {
		return 0, fmt.Errorf("creating the first user: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("committing: %w", err)
	}
	return firstUserCreated, nil
}
