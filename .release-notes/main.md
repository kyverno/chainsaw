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

- Added templating support in resource reference in all collectors
- Added templating support for all string fields in all collectors
- Relaxed `format` field in collectors to support templating

## 🔧 Fixes 🔧

- Added `config` in the list of forbidden binding names
- Fixed resource templating always enabled in `create` operation, regardless of the configured `template` field
- Fixed resource templating always enabled in `patch` operation, regardless of the configured `template` field

## ⭐ Examples ⭐

- Added an example showcasing `x_k8s_server_version` usage
