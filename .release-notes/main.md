# Release notes

Release notes for `TODO`.

<!--
## ⭐ Examples ⭐

## ⛵ Tutorials ⛵

## 📚 Docs 📚

## 🎸 Misc 🎸
-->

## ‼️ Breaking chages ‼️

- Renamed `chainsaw kuttl migrate config` command to  `chainsaw migrate kuttl config`
- Renamed `chainsaw kuttl migrate tests` command to  `chainsaw migrate kuttl tests`

## 💫 New features 💫

- Added a cleanup delay option to add a delay between the time a test ends and the time cleanup starts
- Added support for overriding termination grace period on a per test basis

## 🔧 Fixes 🔧

- Fixed `command` and `script` log output when multiple sections are present
- Fixed delete timeout used instead of cleanup timeout
- Return an error if a test dir is not accessible

## ✨ UI changes ✨

- Use `@setup` and `@cleanup` in the logs to indicate setup/cleanup steps
