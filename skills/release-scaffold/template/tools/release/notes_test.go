package main

import (
	"fmt"
	"strings"
	"testing"
)

// squashPR commits to main the way a squash merge does: the PR's title with
// its number.
func squashPR(t *testing.T, dir string, number int, title string) {
	t.Helper()
	gitIn(t, dir, "commit", "-q", "--allow-empty", "-m", fmt.Sprintf("%s (#%d)", title, number), "-m", "- the squashed body")
}

// mergePR merges a one-commit branch into main with no fast-forward, the
// subject and body as given.
func mergePR(t *testing.T, dir, branch string, message ...string) {
	t.Helper()
	gitIn(t, dir, "checkout", "-q", "-b", branch)
	gitIn(t, dir, "commit", "-q", "--allow-empty", "-m", "work on "+branch)
	gitIn(t, dir, "checkout", "-q", "main")
	args := []string{"merge", "-q", "--no-ff"}
	for _, m := range message {
		args = append(args, "-m", m)
	}
	gitIn(t, dir, append(args, branch)...)
}

// notesFor is the notes of version in dir.
func notesFor(t *testing.T, dir, version string) string {
	t.Helper()
	t.Chdir(dir)
	got, err := notes(version)
	if err != nil {
		t.Fatalf("notes(%q): %v", version, err)
	}
	return got
}

func TestNotesOfTheFirstReleaseListEveryCommit(t *testing.T) {
	dir := newRepo(t)
	squashPR(t, dir, 1, "docs: the first plan")
	squashPR(t, dir, 2, "feat: a window that opens")
	got := notesFor(t, dir, "v0.1.0")
	want := "- feat: a window that opens (#2)\n- docs: the first plan (#1)\n- start\n"
	if got != want {
		t.Errorf("notes of the first release =\n%s\nwant\n%s", got, want)
	}
}

// releasedRepo is a repo released as v0.1.0 after PR #1, with PR #2 in
// v0.2.0-rc.1 and PR #3 in v0.2.0-rc.2, and PR #4 merged since.
func releasedRepo(t *testing.T) string {
	t.Helper()
	dir := newRepo(t)
	squashPR(t, dir, 1, "feat: one")
	gitIn(t, dir, "tag", "v0.1.0")
	squashPR(t, dir, 2, "feat: two")
	gitIn(t, dir, "tag", "v0.2.0-rc.1")
	squashPR(t, dir, 3, "feat: three")
	gitIn(t, dir, "tag", "v0.2.0-rc.2")
	squashPR(t, dir, 4, "feat: four")
	return dir
}

func TestNotesOfAFinalAfterRCsCountFromThePreviousFinal(t *testing.T) {
	got := notesFor(t, releasedRepo(t), "v0.2.0")
	want := "- feat: four (#4)\n- feat: three (#3)\n- feat: two (#2)\n"
	if got != want {
		t.Errorf("notes of v0.2.0 =\n%s\nwant\n%s", got, want)
	}
}

func TestNotesOfAnRCListEverythingNewForItsTarget(t *testing.T) {
	got := notesFor(t, releasedRepo(t), "v0.2.0-rc.3")
	want := "- feat: four (#4)\n- feat: three (#3)\n- feat: two (#2)\n"
	if got != want {
		t.Errorf("notes of v0.2.0-rc.3 =\n%s\nwant\n%s", got, want)
	}
}

func TestNotesRefuseAVersionNotAboveTheLastFinal(t *testing.T) {
	t.Chdir(releasedRepo(t))
	for _, version := range []string{"v0.1.0", "v0.1.0-rc.2", "v0.0.9"} {
		if _, err := notes(version); err == nil || !strings.Contains(err.Error(), "not above the last final release, v0.1.0") {
			t.Errorf("notes(%q) after v0.1.0: err = %v, want a refusal naming v0.1.0", version, err)
		}
	}
}

// Forgejo and Gitea put the PR's title in the merge subject and only a
// Reviewed-on trailer in the body.
func TestNotesOfAForgejoMergeUseTheTitleFromItsSubject(t *testing.T) {
	dir := newRepo(t, "v0.1.0")
	mergePR(t, dir, "docs/home",
		"Merge pull request 'docs: the code lives on Forgejo' (#1) from docs/home into main",
		"Reviewed-on: https://forge.example/o/r/pulls/1")
	got := notesFor(t, dir, "v0.1.1")
	if want := "- docs: the code lives on Forgejo (#1)\n"; got != want {
		t.Errorf("notes of a Forgejo merge = %q, want %q", got, want)
	}
}

// GitHub names the PR in the merge subject and puts its title in the body.
func TestNotesOfAGitHubMergeUseTheTitleFromItsBody(t *testing.T) {
	dir := newRepo(t, "v0.1.0")
	mergePR(t, dir, "pr-7", "Merge pull request #7 from o/pr-7", "feat: a window that opens")
	mergePR(t, dir, "pr-8", "Merge pull request #8 from o/pr-8")
	got := notesFor(t, dir, "v0.2.0")
	if want := "- Merge pull request #8 from o/pr-8\n- feat: a window that opens (#7)\n"; got != want {
		t.Errorf("notes of GitHub merges =\n%s\nwant\n%s", got, want)
	}
}

// A bad version is refused before git runs, so the error names the version
// even outside a repo.
func TestNotesRefuseABadVersion(t *testing.T) {
	t.Chdir(t.TempDir())
	for _, version := range []string{"", "v1", "0.1.0", "v0.1.0-beta.1"} {
		if _, err := notes(version); err == nil || !strings.Contains(err.Error(), "want a release tag") {
			t.Errorf("notes(%q): err = %v, want a bad-version refusal", version, err)
		}
	}
}
