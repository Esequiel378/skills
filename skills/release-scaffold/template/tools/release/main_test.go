package main

import (
	"io"
	"os/exec"
	"slices"
	"strings"
	"testing"
)

// newRepo makes a git repo with one empty commit per tag, tagged in order,
// and returns its directory.
func newRepo(t *testing.T, tags ...string) string {
	t.Helper()
	dir := t.TempDir()
	gitIn(t, dir, "init", "-q", "-b", "main")
	gitIn(t, dir, "commit", "-q", "--allow-empty", "-m", "start")
	for _, tag := range tags {
		gitIn(t, dir, "commit", "-q", "--allow-empty", "-m", tag)
		gitIn(t, dir, "tag", tag)
	}
	return dir
}

func gitIn(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", append([]string{"-c", "user.name=t", "-c", "user.email=t@t", "-c", "commit.gpgsign=false", "-c", "tag.gpgsign=false"}, args...)...)
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, out)
	}
}

func TestReleaseTagsListsTheRepoTags(t *testing.T) {
	t.Chdir(newRepo(t, "v0.1.0", "v0.2.0-rc.1"))
	tags, err := releaseTags()
	if err != nil {
		t.Fatal(err)
	}
	slices.Sort(tags)
	if want := []string{"v0.1.0", "v0.2.0-rc.1"}; !slices.Equal(tags, want) {
		t.Errorf("releaseTags() = %q, want %q", tags, want)
	}
}

func TestReleaseTagsRefusesAShallowClone(t *testing.T) {
	src := newRepo(t, "v0.1.0", "v0.2.0")
	clone := t.TempDir()
	gitIn(t, clone, "clone", "-q", "--depth", "1", "--no-tags", "file://"+src, ".")
	t.Chdir(clone)
	if _, err := releaseTags(); err == nil || !strings.Contains(err.Error(), "shallow") {
		t.Errorf("releaseTags() in a shallow clone: err = %v, want a shallow-clone refusal", err)
	}
}

func TestReleaseTagsRefusesAHeadBehindTheLastRelease(t *testing.T) {
	dir := newRepo(t, "v0.1.0", "v0.2.0")
	gitIn(t, dir, "checkout", "-q", "v0.1.0")
	t.Chdir(dir)
	if _, err := releaseTags(); err == nil || !strings.Contains(err.Error(), "doesn't contain v0.2.0") {
		t.Errorf("releaseTags() at an older commit: err = %v, want a refusal naming v0.2.0", err)
	}
}

func TestReleaseTagsIgnoresOtherTags(t *testing.T) {
	dir := newRepo(t, "v0.1.0")
	gitIn(t, dir, "commit", "-q", "--allow-empty", "-m", "elsewhere")
	gitIn(t, dir, "tag", "v-experiment")
	gitIn(t, dir, "checkout", "-q", "v0.1.0")
	t.Chdir(dir)
	if _, err := releaseTags(); err != nil {
		t.Errorf("releaseTags() with a newer non-release tag: %v", err)
	}
}

func TestRunRejectsBadArguments(t *testing.T) {
	for _, args := range [][]string{nil, {"publish"}, {"version", "-bump", "patch", "extra"}, {"version", "-nope"}, {"notes", "-version", "v0.1.0", "extra"}, {"notes", "-nope"}} {
		if err := run(args, io.Discard); err == nil {
			t.Errorf("run(%q) = nil, want an error", args)
		}
	}
}
