package main

import (
	"strings"
	"testing"
)

func TestNextVersion(t *testing.T) {
	tests := []struct {
		name string
		tags []string
		bump string
		rc   bool
		want string
	}{
		{"no tags: patch is the first release", nil, "patch", false, "v0.1.0"},
		{"no tags: major is the first release too", nil, "major", false, "v0.1.0"},
		{"no tags: an rc of the first release", nil, "minor", true, "v0.1.0-rc.1"},
		{"patch", []string{"v0.1.0", "v0.3.0", "v0.2.0"}, "patch", false, "v0.3.1"},
		{"minor resets patch", []string{"v0.3.2"}, "minor", false, "v0.4.0"},
		{"major resets minor and patch", []string{"v0.3.2"}, "major", false, "v1.0.0"},
		{"the first rc of a target", []string{"v0.3.0"}, "minor", true, "v0.4.0-rc.1"},
		{"the next rc of the same target", []string{"v0.3.0", "v0.4.0-rc.1"}, "minor", true, "v0.4.0-rc.2"},
		{"rc numbers compare as numbers", []string{"v0.3.0", "v0.4.0-rc.9", "v0.4.0-rc.10"}, "minor", true, "v0.4.0-rc.11"},
		{"a final promotes the target's rcs", []string{"v0.3.0", "v0.4.0-rc.2"}, "minor", false, "v0.4.0"},
		{"an rc after its final starts the next target", []string{"v0.3.0", "v0.4.0-rc.2", "v0.4.0"}, "minor", true, "v0.5.0-rc.1"},
		{"a bigger bump leaves a smaller target's rcs behind", []string{"v0.3.0", "v0.4.0-rc.2"}, "major", false, "v1.0.0"},
		{"the first release promotes its rcs", []string{"v0.1.0-rc.1"}, "patch", false, "v0.1.0"},
		{"versions compare as numbers", []string{"v0.9.0", "v0.10.0"}, "patch", false, "v0.10.1"},
		{"other tags are ignored", []string{"v0.3.0", "latest", "v1", "v2.0", "v9.0.0-beta.1", "v9.0.0-rc.x", "x9.0.0", "v09.0.0"}, "patch", false, "v0.3.1"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := nextVersion(tt.tags, tt.bump, tt.rc)
			if err != nil {
				t.Fatalf("nextVersion(%q, %q, %v): %v", tt.tags, tt.bump, tt.rc, err)
			}
			if got != tt.want {
				t.Errorf("nextVersion(%q, %q, %v) = %q, want %q", tt.tags, tt.bump, tt.rc, got, tt.want)
			}
		})
	}
}

func TestNextVersionRefuses(t *testing.T) {
	tests := []struct {
		name string
		tags []string
		bump string
		rc   bool
		want string // in the error
	}{
		{"a patch below a pending rc", []string{"v0.3.0", "v0.4.0-rc.2"}, "patch", false, "v0.3.1 would go backwards: v0.4.0-rc.2 is pending"},
		{"an rc below a pending rc", []string{"v0.3.0", "v0.4.0-rc.2"}, "patch", true, "would go backwards"},
		{"the first release below a pending rc", []string{"v0.2.0-rc.1"}, "patch", false, "v0.1.0 would go backwards"},
		{"an unknown bump", []string{"v0.3.0"}, "huge", false, `bump is "huge"`},
		{"an empty bump", nil, "", false, `bump is ""`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := nextVersion(tt.tags, tt.bump, tt.rc)
			if err == nil || !strings.Contains(err.Error(), tt.want) {
				t.Errorf("nextVersion(%q, %q, %v) = %q, %v; want an error containing %q", tt.tags, tt.bump, tt.rc, got, err, tt.want)
			}
		})
	}
}
