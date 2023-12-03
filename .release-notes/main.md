# Release notes

Release notes for `TODO`.

<!--
## ✨ UI changes ✨

## ⭐ Examples ⭐

## ⛵ Tutorials ⛵

## 📚 Docs 📚
-->

## ‼️ Breaking chages ‼️

- Changed `check` to `expect` in delete `operation`
- Kuttl migration command split into `chainsaw kuttl migrate tests` and `chainsaw kuttl migrate config`
- `timeout` was moved inside operations that support a timeout (all operations but sleep)

## 💫 New features 💫

- Added support to register Chainsaw specific JMESPath functions for use in assertion trees
- Added inline manifest support to `assert` and `error` operations
- Added a `description` field to `Test`, `TestStep`, `Operation`, `Catch` and `Finally`
- Added a command to generate tests documentation
- Added `sleep` operation support in `try`, `catch` and `finally`
- Added support to reference remote files

## 🔧 Fixes 🔧

- Don't show `command` or `script` logs if there's nothing to show
- Fixed `timeouts` always set when migrating kuttl tests

## 🎸 Misc 🎸

- We now have a [dedicated slack channel](https://kubernetes.slack.com/archives/C067LUFL43U)