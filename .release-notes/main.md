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

- Deprecated `TestStep` resource support was removed as well as the related command `chainsaw migrate tests`

## 🔧 Fixes 🔧

- Fixed `wait` operation timeout passing the timeout value to the underlying `kubectl` command
- Fixed `wait` operation condition incorrectly formatted with quotes
- Fixed issue reading files from embedded file systems on windows
