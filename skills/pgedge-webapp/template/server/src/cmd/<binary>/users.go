package main

import (
	"context"
	"errors"
	"fmt"
	"io"

	"<MODULE_PATH>/server/src/internal/auth"
)

// runUserCommand dispatches the requested user-management subcommand.
func runUserCommand(ctx context.Context, store *auth.Store, f *Flags, out io.Writer) error {
	switch {
	case f.AddUser:
		return addUser(ctx, store, f)
	case f.UpdateUser:
		return updateUser(ctx, store, f)
	case f.DeleteUser:
		return store.DeleteUser(ctx, f.Username)
	case f.ListUsers:
		return listUsers(ctx, store, out)
	case f.EnableUser:
		return store.SetEnabled(ctx, f.Username, true)
	case f.DisableUser:
		return store.SetEnabled(ctx, f.Username, false)
	case f.SetSuperuser:
		return store.SetSuperuser(ctx, f.Username, true)
	case f.UnsetSuperuser:
		return store.SetSuperuser(ctx, f.Username, false)
	}
	return errors.New("no user command specified")
}

func addUser(ctx context.Context, store *auth.Store, f *Flags) error {
	if f.Username == "" {
		return errors.New("--username is required")
	}
	pw, err := ResolvePassword(f.Password, f.IsSet("password"), f.PasswordFile)
	if err != nil {
		return err
	}
	if pw == "" {
		return errors.New("--password or --password-file is required")
	}
	return store.CreateUser(ctx, auth.CreateUserParams{
		Username: f.Username, Password: pw,
		FullName: f.FullName, Email: f.Email,
		IsSuperuser: f.SetSuperuser,
	})
}

func updateUser(ctx context.Context, store *auth.Store, f *Flags) error {
	if f.Username == "" {
		return errors.New("--username is required")
	}
	if f.FullName != "" || f.Email != "" {
		if err := store.UpdateProfile(ctx, f.Username, f.FullName, f.Email); err != nil {
			return err
		}
	}
	pw, err := ResolvePassword(f.Password, f.IsSet("password"), f.PasswordFile)
	if err != nil {
		return err
	}
	if pw != "" {
		if err := store.UpdatePassword(ctx, f.Username, pw); err != nil {
			return err
		}
	}
	return nil
}

func listUsers(ctx context.Context, store *auth.Store, out io.Writer) error {
	users, err := store.ListUsers(ctx)
	if err != nil {
		return err
	}
	fmt.Fprintf(out, "%-20s %-30s %-10s %s\n", "USERNAME", "EMAIL", "SUPERUSER", "ENABLED")
	for _, u := range users {
		fmt.Fprintf(out, "%-20s %-30s %-10t %t\n", u.Username, u.Email, u.IsSuperuser, u.Enabled)
	}
	return nil
}
