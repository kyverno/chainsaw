# Release notes

Release notes for `TODO`.

<!--
## ✨ UI changes ✨

## ⭐ Examples ⭐

## ⛵ Tutorials ⛵

-->

## ‼️ Breaking changes ‼️

- Deprecated `TestStep` resource support was removed as well as the related command `chainsaw migrate tests`

## 💫 New features 💫

- Added support to wait with `jsonpath` conditions
- Introduced `v1alpha2` version for configuration, both `v1alpha1` and `v1alpha2` are supported with a conversion system
- Rewrote quick-start guide leveraging chainsaw capabilities using relevant usecases

## 🔧 Fixes 🔧

- Fixed issue with `update` operation failing because resource version not set correctly
- Fixed `wait` operation timeout passing the timeout value to the underlying `kubectl` command
- Fixed `wait` operation condition incorrectly formatted with quotes
- Fixed issue reading files from embedded file systems on Windows
- Fixed incorrect JUnit report format

## 🎸 Misc 🎸

- Switched a release branch worflow

## 📚 Docs 📚

- Added missing cluster doc
