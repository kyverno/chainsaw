# Release notes

Release notes for `TODO`.

<!--
## ‼️ Breaking changes ‼️

## ✨ UI changes ✨

## ⭐ Examples ⭐

## ⛵ Tutorials ⛵
-->

## 💫 New features 💫

- Added multi-cluster support
- Added a new `get` collector
- Added a new `patch` operation
- Added support for `delete` operation in `catch` and `finally` blocks
- Added format support (`json` or `yaml`) in `get` and `events` collectors
- Added support for reading input from stdin in `chainsaw assert` command
- Added programmatic `Configuration` validation
- Added `--report-path` flag and config equivalent to configure the report path

## 🔧 Fixes 🔧

- Fixed a concurrency issue in Kyverno-JSON
- Fixed an issue with timeouts not correctly considered in collectors
- Fixed `--test-dir` flag not parsed correctly when folders were separated by a comma
- Preserve `apiVersion` and `kind` when showing resource diff
- Fixed an incorrect KUTTL `TestAssert` conversion

## 📚 Docs 📚

- Fixed incorrect `delete` operation in docs
- Fixed some typos in the docs

## 🎸 Misc 🎸

- Added docker image labels
- Bumped `ko` and `kind` tools versions
