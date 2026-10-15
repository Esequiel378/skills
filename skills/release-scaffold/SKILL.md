---
name: release-scaffold
description: Use when a project needs a release process, versioned downloads, or a "Release" button — triggers on "/release-scaffold", "add a release process", "set up releases", "ship binaries", "how do we cut a release", or a repo with no release workflow on GitHub or Forgejo.
---

# Release Scaffold

Sets up a release that starts from a button and can't go wrong in the usual
ways. Nobody types a version. A failed run leaves no tag. The number never
goes backwards, and the binary knows its own version.

Most hand-rolled release workflows miss the same things. The version is a
free-text input, so `v0.0.1` can ship after `v1.0.0`. There are no release
candidates. The version logic is bash nobody tests. A failed upload leaves a
draft that blocks the next run. A re-run of the publish job ships a stale
version. The templates here close each of those.

## What the process guarantees

| Guarantee | How |
| --- | --- |
| Only a click releases, only from `main` | `workflow_dispatch` with `bump` (patch/minor/major) and an `rc` checkbox; the first step refuses other refs |
| The version is computed, never typed | `tools/release version`: last *final* tag + bump; `-rc` makes `<target>-rc.N`; a final promotes its rcs; going backwards is refused. Table-tested |
| The tag comes last | Draft at the dispatched commit → upload → publish, which creates the tag. A failed run leaves no tag |
| Never clobbers | Only a 404 counts as a free tag. A draft left by a failed run is deleted first |
| One release at a time | Concurrency group `release`, `cancel-in-progress: false` |
| A re-run can't ship a stale number | "Still the next version" recomputes before publishing |
| A broken commit doesn't ship | The build job runs the same checks CI runs |
| The binary knows its version | Stamped at build time; one built binary is run and must print it |
| Notes on any forge | `tools/release notes`: `main`'s first-parent subjects since the last final tag |

## Files

- [template/tools/release/](template/tools/release/) — the tool: `version` and `notes`, standard library only, with tests. Go 1.24+, whatever the project's language.
- [template/release.github.yml](template/release.github.yml) — publishes with `gh`.
- [template/release.forgejo.yml](template/release.forgejo.yml) — publishes with `curl` against Forgejo's API. Full `https://data.forgejo.org/` action URLs.

## Steps

1. Read the repo. Name its forge, its language, how it builds, and what CI runs.
2. Run `git tag --list 'v*'`. If tags exist that aren't `vX.Y.Z` or `vX.Y.Z-rc.N`, ask the user. The tool ignores them.
3. If the repo has no `v*` tags, tell the user the first release is `v0.1.0` whatever the bump.
4. Copy `template/tools/release/` into the repo as `tools/release/`. The tool needs Go 1.24 or later.
5. If the repo has no root `go.mod`, give `tools/release/` its own: `module <repo>/tools/release` and `go 1.24`.
6. In that case, replace `go run ./tools/release` with `go -C tools/release run .` everywhere: workflow, usage text, doc comment.
7. In that case, point every `go-version-file` at `tools/release/go.mod`, with `cache: false`. Add `go -C tools/release test .` to the build job's checks.
8. If the forge is GitHub, copy `release.github.yml` to `.github/workflows/release.yml`.
9. If the forge is Forgejo, copy `release.forgejo.yml` to `.forgejo/workflows/release.yml`.
10. If the forge is Forgejo, ask the user to confirm the server is v14 or later.
11. If the default branch isn't `main`, rename `main` in the workflow and in the tool's comments.
12. Replace each `# PROJECT:` step in the build job with the project's real commands.
13. Make the checks the ones CI already runs. If the project has no CI, run its test command and any linter it already configures.
14. If the build needs another toolchain, add its setup action. On Forgejo, pin a tag that `https://data.forgejo.org/actions/<name>` has.
15. Build every platform the project ships into `dist/`. If the project ships nothing yet, ask the user which platforms.
16. Name each per-platform file `<app>_<version>_<os>_<arch>`, with `.exe` on Windows. Name a platform-independent file `<app>_<version>.<ext>`. For Go, build with `-trimpath -ldflags "-s -w"` and `CGO_ENABLED=0` unless it needs cgo.
17. Stamp the version into each binary. If the program has no version variable, add one that defaults to `dev`.
18. Make the program print `<app> <version>` for `--version`. `<app>` is the package name; `<version>` keeps its `v`.
19. If a manifest rejects the `v`, write `${VERSION#v}` there.
20. If the language has no build-time flag like Go's `-X`, write the version into a copy of the source or the package manifest during the build.
21. Keep the step that runs one built file and checks it prints `<app> $VERSION`. Run the one the runner can execute.
22. Replace every remaining `APP` with the project's name.
23. Order the build steps so each step's inputs exist before it runs. A test that reads installed dependencies runs after the install step.
24. Add a "Release" section to the README: where the button is, the two inputs, and how to preview the version.
25. If the repo keeps decision records, record the choices in one.

## Verify before you finish

1. Run `go test ./tools/release` (or `go -C tools/release test .`).
2. Run `go run ./tools/release version -bump patch` and show the result.
3. Copy the working tree into a new directory, without build outputs. Run the build job's commands there.
4. If the runner is Linux and a Linux shell is available, run them on Linux.
5. Show the user what this run could not prove: the workflow itself runs only once it is on `main`.
6. Suggest a first dispatch with `rc` ticked. It proves the pipeline without spending a final version.

## When a platform can't build on the runner

Cross-compile when the toolchain allows it. If tests need that platform and no
runner has it, the release can't run them. Say so in the README and in the
decision record. Don't let the gap pass silently.

## Common mistakes

| Mistake | Why it hurts |
| --- | --- |
| Triggering on a pushed tag | The tag exists before the build. A failure leaves an orphan tag to delete by hand |
| `gh release create` without `--draft` | GitHub creates the tag at once, before the upload |
| `forgejo-release` action | It creates the tag before the upload, and reuses one at the same commit |
| `fetch-depth: 1` in the version or release job | The tool refuses a shallow clone: missing tags would look like no releases |
| Version logic in workflow bash | It can't be tested without releasing |
| `--generate-notes` only | GitHub-only, and rc-to-final notes then skip what the rcs shipped |
| No version check on the built binary | A wrong `-X` path silently ships `dev` |
