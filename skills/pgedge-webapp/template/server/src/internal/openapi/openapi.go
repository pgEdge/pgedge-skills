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
	}
}
