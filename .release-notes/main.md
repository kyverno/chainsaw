# Release notes

Release notes for `TODO`.

<!--
## ✨ UI changes ✨

## ⭐ Examples ⭐

## ⛵ Tutorials ⛵

## 📚 Docs 📚

## 🎸 Misc 🎸
-->

## ‼️ Breaking changes ‼️

- Resource templating is now enabled by default

## 💫 New features 💫

- Added `--pause-on-failure` flag to pause when a test failure happens (to ease troubleshooting)
- Improved cleanup management logic, alternating `catch`, `finally` and `@cleanup` per step

## 🔧 Fixes 🔧

- Fixed issue with cluster incorrectly registered
- Force background deletion propagation policy (useful when testing unmanaged `Job`)
