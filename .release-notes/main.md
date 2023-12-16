# Release notes

Release notes for `TODO`.

<!--
## 💫 New features 💫

## ✨ UI changes ✨

## ⭐ Examples ⭐

## ⛵ Tutorials ⛵

## 📚 Docs 📚

## 🎸 Misc 🎸
-->

## ‼️ Breaking changes ‼️

- GitHub action was moved to a dedicated repository https://github.com/kyverno/action-install-chainsaw
- Renamed `chainsaw generate` commands to `chainsaw build`

## 💫 New features 💫

- Added `--readme-file` flag to `chainsaw generate docs` command to customize the name of the generated file
- Added `--catalog` flag to `chainsaw generate docs` command to generate a test catalog
- Added `chainsaw assert` command to evaluate an assertion against a cluster
- Added `chainsaw export schemas` command to copy JSON schemas on disc
- Added `chainsaw create test` command to bootstrap test file
- Added support for labels selector (`--selector` flag) to filter tests to run

## 🔧 Fixes 🔧

- Fixed an invalid error check in `chainsaw docs` command
- Fixed skip delete not applied to test namespace
- Fixed crash due to bindings being null
