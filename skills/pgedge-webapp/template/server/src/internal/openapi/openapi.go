// Package openapi provides a programmatic OpenAPI 3.0.3 specification
// builder for the <PROJECT_NAME> REST API. The same builder powers both
// the runtime HTTP endpoint (/api/v1/openapi.json) and the build-time
// CLI flag (--generate-openapi-spec) used by `make openapi` to refresh
// the docs-tree copy of the spec.
package openapi

import (
	"encoding/json"
)

// M is a shorthand for nested map structures in the spec.
type M = map[string]interface{}

// A is a shorthand for slices in the spec.
type A = []interface{}

// APIVersion is the version reported in the spec's info block. Override
// at build time with -ldflags '-X .../openapi.APIVersion=...' if needed.
var APIVersion = "1.0.0"

// BuildSpec returns the complete OpenAPI 3.0.3 specification as a map.
func BuildSpec() M {
	return M{
		"openapi": "3.0.3",
		"info":    buildInfo(),
		"paths":   buildPaths(),
		"components": M{
			"securitySchemes": buildSecuritySchemes(),
			"schemas":         buildSchemas(),
		},
		"tags": buildTags(),
	}
}

// MarshalJSON returns the spec marshaled as indented JSON.
func MarshalJSON() ([]byte, error) {
	return json.MarshalIndent(BuildSpec(), "", "  ")
}

func buildInfo() M {
	return M{
		"title":       "<PROJECT_NAME> API",
		"description": "REST API for <PROJECT_NAME>.",
		"version":     APIVersion,
		"license":     M{"name": "PostgreSQL License"},
		"contact":     M{"name": "pgEdge", "url": "https://www.pgedge.com"},
	}
}

func buildTags() A {
	return A{
		M{"name": "Health", "description": "Health probe endpoints."},
		M{"name": "Auth", "description": "Session login and logout."},
		M{"name": "User", "description": "Authenticated user information."},
		M{"name": "Users", "description": "User management (superuser only)."},
		M{"name": "OpenAPI", "description": "OpenAPI specification endpoint."},
	}
}

func buildSecuritySchemes() M {
	return M{
		"cookieAuth": M{
			"type":        "apiKey",
			"in":          "cookie",
			"name":        "<COOKIE_NAME>",
			"description": "Session cookie set by POST /api/v1/auth/login.",
		},
	}
}

func ref(name string) M {
	return M{"$ref": "#/components/schemas/" + name}
}

func jsonContent(schema M) M {
	return M{"application/json": M{"schema": schema}}
}

func sessionAuth() A {
	return A{M{"cookieAuth": A{}}}
}

func buildPaths() M {
	return M{
		"/health": M{
			"get": M{
				"tags":        A{"Health"},
				"summary":     "Liveness probe.",
				"description": "Returns 200 OK without authentication; suitable for kubelet liveness/readiness probes.",
				"responses": M{
					"200": M{
						"description": "Server is up.",
						"content":     jsonContent(ref("Health")),
					},
				},
			},
		},
		"/api/v1/auth/login": M{
			"post": M{
				"tags":        A{"Auth"},
				"summary":     "Log in with username and password.",
				"description": "On success sets the session cookie. Subject to a per-IP sliding-window rate limit.",
				"requestBody": M{
					"required": true,
					"content":  jsonContent(ref("LoginRequest")),
				},
				"responses": M{
					"200": M{"description": "Login succeeded.", "content": jsonContent(ref("LoginResponse"))},
					"400": M{"description": "Malformed request body.", "content": jsonContent(ref("ErrorResponse"))},
					"401": M{"description": "Invalid credentials.", "content": jsonContent(ref("ErrorResponse"))},
					"403": M{"description": "Account disabled.", "content": jsonContent(ref("ErrorResponse"))},
					"429": M{"description": "Account locked or rate limited.", "content": jsonContent(ref("ErrorResponse"))},
				},
			},
		},
		"/api/v1/auth/logout": M{
			"post": M{
				"tags":        A{"Auth"},
				"summary":     "Log out the current session.",
				"description": "Clears the session cookie and deletes the matching session row. Idempotent.",
				"responses": M{
					"200": M{"description": "Logout succeeded.", "content": jsonContent(ref("Success"))},
				},
			},
		},
		"/api/v1/user/info": M{
			"get": M{
				"tags":        A{"User"},
				"summary":     "Return the authenticated user's profile.",
				"description": "Requires a valid session cookie.",
				"security":    sessionAuth(),
				"responses": M{
					"200": M{"description": "User info.", "content": jsonContent(ref("UserInfo"))},
					"401": M{"description": "Not authenticated.", "content": jsonContent(ref("UnauthenticatedUserInfo"))},
				},
			},
		},
		"/api/v1/hello": M{
			"get": M{
				"tags":        A{"User"},
				"summary":     "Placeholder protected endpoint.",
				"description": "Demonstrates the authenticated request flow. Replace with your own endpoints.",
				"security":    sessionAuth(),
				"responses": M{
					"200": M{"description": "Hello message.", "content": jsonContent(ref("HelloResponse"))},
					"401": M{"description": "Not authenticated.", "content": jsonContent(ref("ErrorResponse"))},
				},
			},
		},
		"/api/v1/users": M{
			"get": M{
				"tags":        A{"Users"},
				"summary":     "List all users.",
				"description": "Returns all users ordered by username. Requires superuser.",
				"security":    sessionAuth(),
				"responses": M{
					"200": M{"description": "User list.", "content": jsonContent(M{"type": "array", "items": ref("UserListItem")})},
					"401": M{"description": "Not authenticated.", "content": jsonContent(ref("ErrorResponse"))},
					"403": M{"description": "Forbidden — not a superuser.", "content": jsonContent(ref("ErrorResponse"))},
				},
			},
			"post": M{
				"tags":        A{"Users"},
				"summary":     "Create a user.",
				"description": "Creates a new user. Requires superuser.",
				"security":    sessionAuth(),
				"requestBody": M{
					"required": true,
					"content":  jsonContent(ref("CreateUserRequest")),
				},
				"responses": M{
					"201": M{"description": "User created.", "content": jsonContent(ref("CreateUserResponse"))},
					"400": M{"description": "Missing username or password.", "content": jsonContent(ref("ErrorResponse"))},
					"401": M{"description": "Not authenticated.", "content": jsonContent(ref("ErrorResponse"))},
					"403": M{"description": "Forbidden — not a superuser.", "content": jsonContent(ref("ErrorResponse"))},
					"409": M{"description": "Username already exists.", "content": jsonContent(ref("ErrorResponse"))},
				},
			},
		},
		"/api/v1/users/{username}": M{
			"patch": M{
				"tags":        A{"Users"},
				"summary":     "Update a user.",
				"description": "Updates profile and/or flags. All fields are optional. Requires superuser.",
				"security":    sessionAuth(),
				"parameters": A{M{
					"name": "username", "in": "path", "required": true,
					"schema": M{"type": "string"},
				}},
				"requestBody": M{
					"required": true,
					"content":  jsonContent(ref("UpdateUserRequest")),
				},
				"responses": M{
					"200": M{"description": "Update applied.", "content": jsonContent(ref("Success"))},
					"400": M{"description": "Safety rail violation.", "content": jsonContent(ref("ErrorResponse"))},
					"401": M{"description": "Not authenticated.", "content": jsonContent(ref("ErrorResponse"))},
					"403": M{"description": "Forbidden — not a superuser.", "content": jsonContent(ref("ErrorResponse"))},
					"404": M{"description": "User not found.", "content": jsonContent(ref("ErrorResponse"))},
				},
			},
			"delete": M{
				"tags":        A{"Users"},
				"summary":     "Delete a user.",
				"description": "Deletes the user and cascades sessions. Requires superuser.",
				"security":    sessionAuth(),
				"parameters": A{M{
					"name": "username", "in": "path", "required": true,
					"schema": M{"type": "string"},
				}},
				"responses": M{
					"204": M{"description": "Deleted."},
					"400": M{"description": "Cannot delete yourself.", "content": jsonContent(ref("ErrorResponse"))},
					"401": M{"description": "Not authenticated.", "content": jsonContent(ref("ErrorResponse"))},
					"403": M{"description": "Forbidden — not a superuser.", "content": jsonContent(ref("ErrorResponse"))},
					"404": M{"description": "User not found.", "content": jsonContent(ref("ErrorResponse"))},
				},
			},
		},
		"/api/v1/users/{username}/password": M{
			"post": M{
				"tags":        A{"Users"},
				"summary":     "Admin password reset.",
				"description": "Resets a user's password without verifying the old one. Requires superuser.",
				"security":    sessionAuth(),
				"parameters": A{M{
					"name": "username", "in": "path", "required": true,
					"schema": M{"type": "string"},
				}},
				"requestBody": M{
					"required": true,
					"content":  jsonContent(ref("AdminPasswordResetRequest")),
				},
				"responses": M{
					"200": M{"description": "Password updated.", "content": jsonContent(ref("Success"))},
					"400": M{"description": "Empty new_password.", "content": jsonContent(ref("ErrorResponse"))},
					"401": M{"description": "Not authenticated.", "content": jsonContent(ref("ErrorResponse"))},
					"403": M{"description": "Forbidden — not a superuser.", "content": jsonContent(ref("ErrorResponse"))},
					"404": M{"description": "User not found.", "content": jsonContent(ref("ErrorResponse"))},
				},
			},
		},
		"/api/v1/user/password": M{
			"post": M{
				"tags":        A{"User"},
				"summary":     "Self-service password change.",
				"description": "Changes the caller's own password after verifying the current one.",
				"security":    sessionAuth(),
				"requestBody": M{
					"required": true,
					"content":  jsonContent(ref("SelfPasswordChangeRequest")),
				},
				"responses": M{
					"200": M{"description": "Password changed.", "content": jsonContent(ref("Success"))},
					"400": M{"description": "Empty new_password.", "content": jsonContent(ref("ErrorResponse"))},
					"401": M{"description": "Not authenticated or wrong current password.", "content": jsonContent(ref("ErrorResponse"))},
				},
			},
		},
		"/api/v1/openapi.json": M{
			"get": M{
				"tags":        A{"OpenAPI"},
				"summary":     "Return this OpenAPI specification as JSON.",
				"description": "Public endpoint; no authentication required.",
				"responses": M{
					"200": M{
						"description": "OpenAPI 3.0.3 specification.",
						"content":     M{"application/json": M{"schema": M{"type": "object"}}},
					},
				},
			},
		},
	}
}

func buildSchemas() M {
	return M{
		"Health": M{
			"type":     "object",
			"required": A{"status"},
			"properties": M{
				"status": M{"type": "string", "example": "ok"},
			},
		},
		"LoginRequest": M{
			"type":     "object",
			"required": A{"username", "password"},
			"properties": M{
				"username": M{"type": "string", "example": "admin"},
				"password": M{"type": "string", "format": "password"},
			},
		},
		"LoginResponse": M{
			"type":     "object",
			"required": A{"success"},
			"properties": M{
				"success":    M{"type": "boolean"},
				"expires_at": M{"type": "string", "format": "date-time"},
				"message":    M{"type": "string"},
			},
		},
		"Success": M{
			"type":     "object",
			"required": A{"success"},
			"properties": M{
				"success": M{"type": "boolean"},
			},
		},
		"UserInfo": M{
			"type":     "object",
			"required": A{"authenticated", "username"},
			"properties": M{
				"authenticated": M{"type": "boolean", "example": true},
				"username":      M{"type": "string"},
				"is_superuser":  M{"type": "boolean"},
			},
		},
		"UnauthenticatedUserInfo": M{
			"type":     "object",
			"required": A{"authenticated"},
			"properties": M{
				"authenticated": M{"type": "boolean", "example": false},
			},
		},
		"HelloResponse": M{
			"type":     "object",
			"required": A{"message", "username"},
			"properties": M{
				"message":  M{"type": "string", "example": "Hello, alice"},
				"username": M{"type": "string", "example": "alice"},
			},
		},
		"ErrorResponse": M{
			"type":     "object",
			"required": A{"error"},
			"properties": M{
				"error":   M{"type": "string"},
				"message": M{"type": "string"},
			},
		},
		"UserListItem": M{
			"type":     "object",
			"required": A{"username", "full_name", "email", "is_superuser", "enabled"},
			"properties": M{
				"username":     M{"type": "string"},
				"full_name":    M{"type": "string"},
				"email":        M{"type": "string"},
				"is_superuser": M{"type": "boolean"},
				"enabled":      M{"type": "boolean"},
			},
		},
		"CreateUserRequest": M{
			"type":     "object",
			"required": A{"username", "password"},
			"properties": M{
				"username":     M{"type": "string", "example": "bob"},
				"password":     M{"type": "string", "format": "password"},
				"full_name":    M{"type": "string"},
				"email":        M{"type": "string", "format": "email"},
				"is_superuser": M{"type": "boolean"},
			},
		},
		"CreateUserResponse": M{
			"type":     "object",
			"required": A{"success", "username"},
			"properties": M{
				"success":  M{"type": "boolean"},
				"username": M{"type": "string"},
			},
		},
		"UpdateUserRequest": M{
			"type": "object",
			"properties": M{
				"full_name":    M{"type": "string"},
				"email":        M{"type": "string", "format": "email"},
				"is_superuser": M{"type": "boolean"},
				"enabled":      M{"type": "boolean"},
			},
		},
		"AdminPasswordResetRequest": M{
			"type":     "object",
			"required": A{"new_password"},
			"properties": M{
				"new_password": M{"type": "string", "format": "password"},
			},
		},
		"SelfPasswordChangeRequest": M{
			"type":     "object",
			"required": A{"current_password", "new_password"},
			"properties": M{
				"current_password": M{"type": "string", "format": "password"},
				"new_password":     M{"type": "string", "format": "password"},
			},
		},
	}
}
