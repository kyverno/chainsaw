# Release notes

Release notes for `TODO`.

<!--
## ‼️ Breaking changes ‼️

## ✨ UI changes ✨

## ⭐ Examples ⭐

## ⛵ Tutorials ⛵
-->

## 💫 New features 💫

- Allowed passing test folders by args (`chainsaw test ./folder` instead of `chainsaw test --test-dir ./folder`)
- Added new binding `$namespace` containing the test namespace name
- Added new `describe` collector to invoke `kubectl describe ...`

## 🔧 Fixes 🔧

- Fixed an issue when running `chainsaw migrate kuttl tests` twice on the same folder
- Fixed an issue with `chainsaw migrate kuttl tests` potentially incorrect when a `TestStep` is present
- Fixed an incorrect `Expectation` API docs
- Fixed an incorrect `FileOrCheck` API docs

## 📚 Docs 📚

- Improved API docs for types coming from Kyverno-JSON

## 🎸 Misc 🎸

- Added an issue template for general question
- Added issue templates config with blank issues and various links
