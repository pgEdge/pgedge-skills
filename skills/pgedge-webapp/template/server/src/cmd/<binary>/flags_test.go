package main

import "testing"

func TestParseFlags_Defaults(t *testing.T) {
	f, err := ParseFlags(nil, "/etc/x.yaml")
	if err != nil {
		t.Fatal(err)
	}
	if f.ConfigFile != "/etc/x.yaml" {
		t.Errorf("default config = %q", f.ConfigFile)
	}
	if f.HasUserCommand() {
		t.Error("no user command should be detected")
	}
}

func TestParseFlags_CLIFlagsOverlay(t *testing.T) {
	f, err := ParseFlags([]string{"-addr", ":7777", "-tls"}, "")
	if err != nil {
		t.Fatal(err)
	}
	cli := f.ToCLIFlags()
	if !cli.HTTPAddrSet || cli.HTTPAddr != ":7777" {
		t.Errorf("addr overlay = %+v", cli)
	}
	if !cli.TLSEnabledSet || !cli.TLSEnabled {
		t.Errorf("tls overlay = %+v", cli)
	}
}

func TestHasUserCommand(t *testing.T) {
	f, _ := ParseFlags([]string{"-add-user"}, "")
	if !f.HasUserCommand() {
		t.Error("expected add-user detected")
	}
}
