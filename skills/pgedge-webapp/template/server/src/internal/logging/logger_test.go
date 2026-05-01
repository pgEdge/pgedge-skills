package logging

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestNew_JSON(t *testing.T) {
	var buf bytes.Buffer
	log := New(&buf, false)
	log.Info("hello", "key", "value")

	var got map[string]any
	if err := json.Unmarshal(buf.Bytes(), &got); err != nil {
		t.Fatalf("output is not JSON: %v\n%s", err, buf.String())
	}
	if got["msg"] != "hello" || got["key"] != "value" {
		t.Errorf("missing fields: %+v", got)
	}
}

func TestNew_Debug(t *testing.T) {
	var buf bytes.Buffer
	log := New(&buf, true)
	log.Debug("dbg", "x", 1)
	if !strings.Contains(buf.String(), "dbg") {
		t.Errorf("debug message not emitted: %q", buf.String())
	}
	if !strings.Contains(buf.String(), "x=1") {
		t.Errorf("debug uses text handler: %q", buf.String())
	}
}

func TestNew_DefaultLevelIsInfo(t *testing.T) {
	var buf bytes.Buffer
	log := New(&buf, false)
	log.Debug("hidden")
	if buf.Len() != 0 {
		t.Errorf("debug should be suppressed at info level: %q", buf.String())
	}
}
