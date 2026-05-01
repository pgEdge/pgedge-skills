package openapi

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestBuildSpec_TopLevelKeys(t *testing.T) {
	spec := BuildSpec()
	for _, k := range []string{"openapi", "info", "paths", "components", "tags"} {
		if _, ok := spec[k]; !ok {
			t.Errorf("missing top-level key %q", k)
		}
	}
	if v, _ := spec["openapi"].(string); !strings.HasPrefix(v, "3.0") {
		t.Errorf("openapi version = %q, want 3.0.x", v)
	}
}

func TestBuildSpec_AllExpectedPaths(t *testing.T) {
	spec := BuildSpec()
	paths, _ := spec["paths"].(M)
	want := []string{
		"/health",
		"/api/v1/auth/login",
		"/api/v1/auth/logout",
		"/api/v1/user/info",
		"/api/v1/hello",
		"/api/v1/openapi.json",
	}
	for _, p := range want {
		if _, ok := paths[p]; !ok {
			t.Errorf("missing path %q", p)
		}
	}
}

func TestMarshalJSON_RoundTrip(t *testing.T) {
	b, err := MarshalJSON()
	if err != nil {
		t.Fatalf("MarshalJSON: %v", err)
	}
	var got map[string]interface{}
	if err := json.Unmarshal(b, &got); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	if got["openapi"] == nil {
		t.Error("round-tripped spec missing openapi key")
	}
}

func TestBuildSpec_CookieSecurityScheme(t *testing.T) {
	spec := BuildSpec()
	components, _ := spec["components"].(M)
	schemes, _ := components["securitySchemes"].(M)
	cookie, ok := schemes["cookieAuth"].(M)
	if !ok {
		t.Fatal("missing cookieAuth security scheme")
	}
	if cookie["in"] != "cookie" {
		t.Errorf("expected in=cookie, got %v", cookie["in"])
	}
	if cookie["type"] != "apiKey" {
		t.Errorf("expected type=apiKey, got %v", cookie["type"])
	}
}
