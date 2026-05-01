package main

import (
	"flag"

	"<MODULE_PATH>/server/src/internal/config"
)

// Flags holds parsed command-line flags.
type Flags struct {
	fs *flag.FlagSet

	ConfigFile        string
	HTTPAddr          string
	TLSEnabled        bool
	CertFile          string
	KeyFile           string
	ChainFile         string
	DataDir           string
	Debug             bool
	AdminPasswordFile string
	Init              bool

	AddUser, DeleteUser, UpdateUser, ListUsers, EnableUser, DisableUser bool
	SetSuperuser, UnsetSuperuser                                        bool
	Username, Password, PasswordFile, FullName, Email                   string
}

// ParseFlags parses argv (excluding the program name) into a Flags struct.
func ParseFlags(args []string, defaultConfigPath string) (*Flags, error) {
	f := &Flags{}
	fs := flag.NewFlagSet("<BINARY_NAME>", flag.ContinueOnError)
	f.fs = fs

	fs.StringVar(&f.ConfigFile, "config", defaultConfigPath, "Path to config file")
	fs.StringVar(&f.HTTPAddr, "addr", "", "HTTP listen address (overrides config)")
	fs.BoolVar(&f.TLSEnabled, "tls", false, "Enable TLS")
	fs.StringVar(&f.CertFile, "cert", "", "TLS certificate file")
	fs.StringVar(&f.KeyFile, "key", "", "TLS key file")
	fs.StringVar(&f.ChainFile, "chain", "", "TLS chain file (optional)")
	fs.StringVar(&f.DataDir, "data-dir", "", "Data directory (sqlite db lives here)")
	fs.BoolVar(&f.Debug, "debug", false, "Verbose request logging")
	fs.StringVar(&f.AdminPasswordFile, "admin-password-file", "",
		"Path to file containing the seed admin password (overrides built-in seed)")
	fs.BoolVar(&f.Init, "init", false, "One-shot: create data dir, init schema, seed admin")

	fs.BoolVar(&f.AddUser, "add-user", false, "")
	fs.BoolVar(&f.DeleteUser, "delete-user", false, "")
	fs.BoolVar(&f.UpdateUser, "update-user", false, "")
	fs.BoolVar(&f.ListUsers, "list-users", false, "")
	fs.BoolVar(&f.EnableUser, "enable-user", false, "")
	fs.BoolVar(&f.DisableUser, "disable-user", false, "")
	fs.BoolVar(&f.SetSuperuser, "set-superuser", false, "")
	fs.BoolVar(&f.UnsetSuperuser, "unset-superuser", false, "")
	fs.StringVar(&f.Username, "username", "", "")
	fs.StringVar(&f.Password, "password", "", "")
	fs.StringVar(&f.PasswordFile, "password-file", "", "")
	fs.StringVar(&f.FullName, "full-name", "", "")
	fs.StringVar(&f.Email, "email", "", "")

	if err := fs.Parse(args); err != nil {
		return nil, err
	}
	return f, nil
}

// IsSet returns whether the named flag was explicitly provided.
func (f *Flags) IsSet(name string) bool {
	set := false
	f.fs.Visit(func(fl *flag.Flag) {
		if fl.Name == name {
			set = true
		}
	})
	return set
}

// HasUserCommand reports whether any user-management subcommand was given.
func (f *Flags) HasUserCommand() bool {
	return f.AddUser || f.DeleteUser || f.UpdateUser || f.ListUsers ||
		f.EnableUser || f.DisableUser || f.SetSuperuser || f.UnsetSuperuser
}

// ToCLIFlags converts Flags into the config.CLIFlags overlay.
func (f *Flags) ToCLIFlags() config.CLIFlags {
	return config.CLIFlags{
		ConfigFile:    f.ConfigFile,
		ConfigFileSet: f.IsSet("config"),
		HTTPAddr:      f.HTTPAddr, HTTPAddrSet: f.IsSet("addr"),
		TLSEnabled: f.TLSEnabled, TLSEnabledSet: f.IsSet("tls"),
		TLSCertFile: f.CertFile, TLSCertSet: f.IsSet("cert"),
		TLSKeyFile: f.KeyFile, TLSKeySet: f.IsSet("key"),
		TLSChainFile: f.ChainFile, TLSChainSet: f.IsSet("chain"),
		DataDir: f.DataDir, DataDirSet: f.IsSet("data-dir"),
	}
}
