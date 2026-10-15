package main

import (
	"fmt"
	"regexp"
	"strings"
)

// forgejoMergePattern matches the subject of a Forgejo or Gitea merge
// commit, which carries the PR's title and number; the body holds only
// trailers.
var forgejoMergePattern = regexp.MustCompile(`^Merge pull request '(.+)' \(#([0-9]+)\) from .+ into .+$`)

// githubMergePattern matches the subject of a GitHub merge commit and
// captures its PR number; GitHub puts the PR's title in the body.
var githubMergePattern = regexp.MustCompile(`^Merge pull request #([0-9]+) from `)

// notes is the Markdown release notes of version at HEAD: one line per
// commit of main's first-parent history since the last final release,
// newest first, so an rc lists everything new for its target and the first
// release lists everything. A squash merge is its subject, which already
// ends in (#N); a merge commit is its PR's title and number. The tags alone
// set the range: version only guards against notes for a version already
// behind the last final.
func notes(version string) (string, error) {
	target, _, ok := parseTag(version)
	if !ok {
		return "", badVersion("version", version)
	}
	tags, err := releaseTags()
	if err != nil {
		return "", err
	}
	since := "HEAD"
	if last, released := lastFinal(tags); released {
		if !last.less(target) {
			return "", fmt.Errorf("%s is not above the last final release, %s", version, last)
		}
		since = last.String() + "..HEAD"
	}
	// Each field ends with a NUL, which git refuses in a commit message.
	out, err := git("log", "--first-parent", "-z", "--format=%s%x00%b", since, "--")
	if err != nil {
		return "", err
	}
	fields := strings.Split(out, "\x00")
	var b strings.Builder
	for i := 0; i+1 < len(fields); i += 2 {
		subject, body := strings.TrimLeft(fields[i], "\n"), fields[i+1]
		if m := forgejoMergePattern.FindStringSubmatch(subject); m != nil {
			subject = m[1] + " (#" + m[2] + ")"
		} else if m := githubMergePattern.FindStringSubmatch(subject); m != nil {
			if title := firstLine(body); title != "" {
				subject = title + " (#" + m[1] + ")"
			}
		}
		fmt.Fprintf(&b, "- %s\n", subject)
	}
	return b.String(), nil
}

// firstLine is the first non-blank line of s, trimmed.
func firstLine(s string) string {
	for line := range strings.Lines(s) {
		if line = strings.TrimSpace(line); line != "" {
			return line
		}
	}
	return ""
}
