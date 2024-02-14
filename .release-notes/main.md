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
- Added `namespaceTemplate` support in configuration and test to allow customizing the test namespace before creation
- Added `template` support in configuration, test, step and operations to enable resource template processing
- Improved KUTTL tests migration files rewrite logic, keeping the same file names
- Added `--template` flag to enable/disable resource templating support from the command line

## 🔧 Fixes 🔧

- Fixed scripts and collectors not working in docker container
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
- Added makefile target to build a docker image locally with [ko](https://ko.build/)
