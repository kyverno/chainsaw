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
- Added diff output in `assert` when an assertion fails
- Added new binding `$namespace` containing the test namespace name
- Added new `describe` collector to invoke `kubectl describe ...`
- Added `modifiers` support in `apply` and `create` operations to modify resources before they are submitted to the cluster
- Added `namespaceModifiers` support in `Test` and `Configuration` to modify test namespaces before they are submitted to the cluster
- Improved KUTTL tests migration files rewrite logic, keeping the same file names

## 🔧 Fixes 🔧

- Fixed an issue with `List` objects not loaded correctly
- Fixed an issue when running `chainsaw migrate kuttl tests` twice on the same folder
- Fixed an issue with `chainsaw migrate kuttl tests` potentially incorrect when a `TestStep` is present
- Fixed an incorrect `Expectation` API docs
- Fixed an incorrect `FileOrCheck` API docs

## 📚 Docs 📚

- Improved API docs for types coming from Kyverno-JSON

## 🎸 Misc 🎸

- Added an issue template for general question
- Added issue templates config with blank issues and various links
