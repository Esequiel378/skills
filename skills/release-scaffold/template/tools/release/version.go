package main

import (
	"fmt"
	"regexp"
	"strconv"
)

// A release tag is a final vX.Y.Z or a release candidate vX.Y.Z-rc.N.
// Leading zeros aren't semver, so they don't match.
var tagPattern = regexp.MustCompile(`^v(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)(?:-rc\.([1-9][0-9]*))?$`)

// semver is a version without its rc number.
type semver struct{ major, minor, patch int }

func (v semver) String() string { return fmt.Sprintf("v%d.%d.%d", v.major, v.minor, v.patch) }

func (v semver) less(w semver) bool {
	if v.major != w.major {
		return v.major < w.major
	}
	if v.minor != w.minor {
		return v.minor < w.minor
	}
	return v.patch < w.patch
}

// parseTag reads a release tag. rcNumber is 0 for a final; ok is false for
// any other tag, which the release ignores.
func parseTag(tag string) (v semver, rcNumber int, ok bool) {
	m := tagPattern.FindStringSubmatch(tag)
	if m == nil {
		return semver{}, 0, false
	}
	var n [4]int
	for i, s := range m[1:] {
		if s == "" {
			continue // a final has no rc
		}
		var err error
		if n[i], err = strconv.Atoi(s); err != nil {
			return semver{}, 0, false // too big for an int
		}
	}
	return semver{n[0], n[1], n[2]}, n[3], true
}

// badVersion is the error for a version flag or argument that isn't a
// release tag.
func badVersion(what, version string) error {
	return fmt.Errorf("%s is %q, want a release tag, vX.Y.Z or vX.Y.Z-rc.N", what, version)
}

// lastFinal is the highest final release among tags; ok is false when
// nothing has been released yet.
func lastFinal(tags []string) (last semver, ok bool) {
	for _, tag := range tags {
		if v, rcNumber, isRelease := parseTag(tag); isRelease && rcNumber == 0 && (!ok || last.less(v)) {
			last, ok = v, true
		}
	}
	return last, ok
}

// nextVersion is the version a release run makes from the repo's tags.
// The target is the highest final plus bump, or v0.1.0
// when there is no final yet. An rc run makes the target's next rc; a final
// run makes the target, promoting its rcs. A target below a pending rc would
// go backwards, so it is refused.
func nextVersion(tags []string, bump string, makeRC bool) (string, error) {
	last, released := lastFinal(tags)
	var target semver
	switch bump {
	case "patch":
		target = semver{last.major, last.minor, last.patch + 1}
	case "minor":
		target = semver{last.major, last.minor + 1, 0}
	case "major":
		target = semver{last.major + 1, 0, 0}
	default:
		return "", fmt.Errorf("bump is %q, want patch, minor or major", bump)
	}
	if !released {
		target = semver{0, 1, 0} // the first release, whatever the bump
	}

	newestRC := 0
	for _, tag := range tags {
		v, rcNumber, ok := parseTag(tag)
		if !ok || rcNumber == 0 || (released && !last.less(v)) {
			continue // not a pending rc: another tag, a final, or an rc of a released version
		}
		if target.less(v) {
			return "", fmt.Errorf("%s would go backwards: %s is pending; pick a bigger bump", target, tag)
		}
		if v == target && rcNumber > newestRC {
			newestRC = rcNumber
		}
	}
	if makeRC {
		return fmt.Sprintf("%s-rc.%d", target, newestRC+1), nil
	}
	return target.String(), nil
}
