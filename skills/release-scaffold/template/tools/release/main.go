// Command release holds the release workflow's logic, so it is tested under
// go test and runs the same on any forge or locally.
//
//	go run ./tools/release version -bump <patch|minor|major> [-rc]
//
// prints the version the next release of HEAD makes from the repo's tags, or
// exits non-zero when the run must be refused.
//
//	go run ./tools/release notes -version <vX.Y.Z[-rc.N]>
//
// prints that release's notes as Markdown: what main has merged since the
// last final release.
package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

const usage = "usage: go run ./tools/release version -bump <patch|minor|major> [-rc] | notes -version <vX.Y.Z[-rc.N]>"

func main() {
	if err := run(os.Args[1:], os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, "release:", err)
		os.Exit(1)
	}
}

// run runs the command args, printing its result to stdout.
func run(args []string, stdout io.Writer) error {
	if len(args) == 0 {
		return errors.New(usage)
	}
	switch args[0] {
	case "version":
		fs := flag.NewFlagSet("version", flag.ContinueOnError)
		bump := fs.String("bump", "", "patch, minor or major")
		makeRC := fs.Bool("rc", false, "make a release candidate")
		if err := parse(fs, args[1:]); err != nil {
			return err
		}
		tags, err := releaseTags()
		if err != nil {
			return err
		}
		v, err := nextVersion(tags, *bump, *makeRC)
		if err != nil {
			return err
		}
		fmt.Fprintln(stdout, v)
		return nil
	case "notes":
		fs := flag.NewFlagSet("notes", flag.ContinueOnError)
		version := fs.String("version", "", "the version being released")
		if err := parse(fs, args[1:]); err != nil {
			return err
		}
		text, err := notes(*version)
		if err != nil {
			return err
		}
		fmt.Fprint(stdout, text)
		return nil
	default:
		return fmt.Errorf("unknown command %q; %s", args[0], usage)
	}
}

// parse parses a subcommand's flags, which take no arguments.
func parse(fs *flag.FlagSet, args []string) error {
	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() > 0 {
		return fmt.Errorf("%s: unexpected argument %q", fs.Name(), fs.Arg(0))
	}
	return nil
}

// releaseTags lists the repo's tags for a release of HEAD. It refuses a
// shallow clone, whose missing tags would look like a repo that never
// released, and a HEAD that doesn't contain every release so far, which
// would tag an older commit after a newer one. Every release is made from
// main, so main's tip always contains them.
func releaseTags() ([]string, error) {
	shallow, err := git("rev-parse", "--is-shallow-repository")
	if err != nil {
		return nil, err
	}
	if shallow == "true" {
		return nil, errors.New("this is a shallow clone, so its tags are incomplete; fetch the full history")
	}
	out, err := git("tag", "--list", "v*")
	if err != nil {
		return nil, err
	}
	tags := strings.Fields(out)
	for _, tag := range tags {
		if _, _, ok := parseTag(tag); !ok {
			continue
		}
		_, err := git("merge-base", "--is-ancestor", tag, "HEAD")
		var exitErr *exec.ExitError
		switch {
		case errors.As(err, &exitErr) && exitErr.ExitCode() == 1:
			return nil, fmt.Errorf("HEAD doesn't contain %s, an earlier release: release a newer commit", tag)
		case err != nil:
			return nil, err
		}
	}
	return tags, nil
}

// git runs a git command and returns its trimmed output; an error carries
// git's own message.
func git(args ...string) (string, error) {
	out, err := exec.Command("git", args...).Output()
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) && len(exitErr.Stderr) > 0 {
			return "", fmt.Errorf("git %s: %w: %s", args[0], err, bytes.TrimSpace(exitErr.Stderr))
		}
		return "", fmt.Errorf("git %s: %w", args[0], err)
	}
	return strings.TrimSpace(string(out)), nil
}
