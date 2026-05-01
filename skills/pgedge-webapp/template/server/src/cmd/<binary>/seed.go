package main

// SeededAdminPasswordHash holds the bcrypt hash of the initial admin password.
// The scaffolding tool replaces <SEEDED_ADMIN_PASSWORD_HASH> with the actual
// hash generated from the password supplied to the skill at scaffold time.
//
// Rotate after first login. Subsequent password changes go through
// `<BINARY_NAME> -update-user -username admin -password-file new.txt`.
const SeededAdminPasswordHash = "<SEEDED_ADMIN_PASSWORD_HASH>"

// SeededAdminUsername is the username for the seed admin account.
const SeededAdminUsername = "admin"
