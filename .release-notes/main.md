# Release notes

Release notes for `TODO`.

<!--
## ‼️ Breaking changes ‼️

## ✨ UI changes ✨

## ⛵ Tutorials ⛵

## 📚 Docs 📚

## 🎸 Misc 🎸
-->

## 💫 New features 💫

- Added templating support in `bindings` and `outputs` names
- Added templating support in `env` names
- Added templating support in resource reference in all collectors
- Added templating support for all string fields in all collectors
- Removed bindings forbidden names (be careful if you override a built-in binding name)
- Added `match` support in `command` and `script` operations outputs
- Relaxed `format` field validation in collectors to support templating
- Relaxed bindings `name` field validation to support templating
- Relaxed env `name` field validation to support templating

## 🔧 Fixes 🔧

- Fixed a resource templating issue in non-resource assertions
- Fixed diff not templated in case of `assert` failure
- Fixed resource templating always enabled in `create` operation, regardless of the configured `template` field
- Fixed resource templating always enabled in `patch` operation, regardless of the configured `template` field
- Added missing operations to the `build docs` command template
- Added test-level catch statements to the `build docs` command template
- Added binding `name` validation markers
- Fixed `build docs` command for files with multiple tests
- Relaxed step file discovery regular expression to allow names with one or more digit prefixes

## ⭐ Examples ⭐

- Added an example showcasing `x_k8s_server_version` usage
