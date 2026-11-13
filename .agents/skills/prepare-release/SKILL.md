---
name: prepare-release
description: Cuts a Gomander release end to end, from the version bump PR to a draft release with written notes and every binary uploaded. Use when the user asks to prepare, cut or ship a release.
disable-model-invocation: true
---

# Prepare release

Invoking this skill authorises opening the bump PR, merging it once green and editing the draft release. It does not authorise publishing the release: that stays with the user.

## How a release happens here

- The version is the single constant `CurrentRelease` in `internal/releases/release.go`. Nothing else is bumped by hand; `scripts/create-deb.sh` reads it from there.
- Merging a change to that file on `main` fires `tag-on-version-bump.yml`, which pushes the tag `vX.Y.Z`.
- The tag fires `release.yml`: `create-release` makes a draft with `--generate-notes`, then the darwin, windows and linux (amd64, arm64) jobs upload into it.
- A complete draft has 12 assets: `gomander-darwin-{amd64,arm64}.{dmg,app.zip}`, `gomander-windows-{amd64,arm64}-{portable,installer}.exe`, `gomander-linux-{amd64,arm64}` and `gomander_X.Y.Z_{amd64,arm64}.deb`.

## Steps

1. **Changes since the last release.** Fetch, take the latest published tag (`gh release list`), then the PRs merged into `main` after it (`git log --first-parent <tag>..origin/main`, then each PR). Classify each PR by its title's conventional type, not by the CodeRabbit summary in its body, which files fixes under "New Features". If nothing user-facing landed, say so and stop.
2. **Propose the bump and wait for the user's answer.** Give the version, the reason and the PRs behind it, and stop until they confirm or change it.
   - major: something stops working for an existing user, such as a dropped OS version (v2 dropped macOS < 12), a third-party API (`cmd/gomander/thirdpartyserver/openapi.yaml`) change that breaks callers, or a project export file that an older or newer version can no longer import.
   - minor: any `feat`.
   - patch: only fixes and internal changes.
3. **Bump PR.** Branch `release-vX.Y.Z` from `origin/main`, change only `CurrentRelease`, commit `chore(release): bump version to vX.Y.Z`. Open it with `.github/pull_request_template.md`, as #354 did: what ships, one line per PR, and in *Additional Notes* that merging it is the release trigger. It is not a draft.
4. **Wait for every check to pass, then merge** with a merge commit. A red check stops the flow: report it, don't fix it here.
5. **Confirm the tag.** Within a few minutes `vX.Y.Z` should exist and a `Release` run should start (`gh run list --workflow release.yml`). A tag with no run means `RELEASE_PERSONAL_ACCESS_TOKEN` expired and the tag was pushed with `GITHUB_TOKEN`, which triggers nothing: tell the user to renew it (steps in `tag-on-version-bump.yml`) and stop.
6. **Write the notes while the binaries build.** Wait until the draft exists (`create-release` finishes in seconds), then read the bodies of the last two published releases and follow their structure and tone:
   - `## TLDR`: what a user notices, in a short paragraph or two, then whether rolling back is safe: what the previous version does with anything this one writes, such as columns from `migrations/` or a new export file version.
   - `## ⚠️ Why this is a major version`, only for a major.
   - `## ✨ New Features`, `## 🐛 Fixes`, `## 🧹 Changes & Removals`, `## 🏗️ Under the hood`: bold lead, em dash, user-level detail, PR numbers at the end. Leave out a section with nothing in it.
   - The draft's generated notes, unchanged (they include the bump PR), inside the `<details><summary><h2>Full changelog</h2></summary>` block.
   - The closing `## Important reminder for macOS users` section linking the README.

   Write them in English and set them with `gh release edit vX.Y.Z --notes-file <file>`, keeping the draft a draft.
7. **Wait for the run to finish**, then check the draft has the 12 assets. A failed job means a broken release: a tag can't be rebuilt, so the fix ships as the next patch (v2.0.0 was never published; v2.0.1 replaced it). Report which job failed and stop.
8. **Tell the user** it is ready: the draft's URL, the asset count and a summary of the notes, so they can review and publish it.
