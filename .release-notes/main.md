# Release notes

Release notes for `TODO`.

<!--
## ‼️ Breaking changes ‼️

## ✨ UI changes ✨

## ⭐ Examples ⭐

## ⛵ Tutorials ⛵

## 📚 Docs 📚

## 🎸 Misc 🎸
-->

## 💫 New features 💫

- Added `chainsaw renovate config` to upgrade configuration manifests to `v1alpha2`
- Changed `cluster` in test step and actions to a pointer
    - `null` means inherit the current cluster
    - an empty string means the default cluster
- Made default timeouts part of the schemas
- Added support to fail fast at the test level (test will be skipped if a previous error was already reported)
- Delay operation loading until operation execution

## 🔧 Fixes 🔧

- Fixed incorrect bindings when operations are run against different clusters
- Fixed operation info built-in bindings not registered before other bindings
