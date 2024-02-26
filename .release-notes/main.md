# Release notes

Release notes for `TODO`.

<!--
## ‼️ Breaking changes ‼️

## ✨ UI changes ✨

## ⭐ Examples ⭐

## ⛵ Tutorials ⛵
-->

## 💫 New features 💫

- Added support for custom bindings (local variables) at the test, step and operation level
- Added support for custom environment variables in `script` and `command` operations
- Added support for loading resources from a file in `chainsaw assert` command
- Added new `wait` collector to invoke `kubectl wait ...` 

## 🔧 Fixes 🔧

- Fixed report path configuration not being used when saving a report
- Prevent double extension in the report output file
- Fixed a potential crash related to reports
- Fixed potentially missing operations in reports

## 🎸 Misc 🎸

- Temporary `KUBECONFIG` files are now created in `/tmp` folder
- Stopped overriding the `PATH` environment variable

## 📚 Docs 📚

- Added chainsaw capabilities in the quick start guide
