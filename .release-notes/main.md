# Release notes

Release notes for `TODO`.

<!--
## ‼️ Breaking changes ‼️

## ✨ UI changes ✨

## ⭐ Examples ⭐

## ⛵ Tutorials ⛵

## 📚 Docs 📚
-->

## 💫 New features 💫

- Support both `.yaml` and `.yml` extensions for chainsaw test files
- Support specifying `apiVersion` and `kind` in collectors (working with resources should be considered deprecated)
- Support all namespaces in collectors (by specifying `namespace: '*'`)
- Added `outputs` support in `script` and `command` operations
- Added support to specify `catch` blocks at the configuration and test levels
- Support waiting `all` in `wait` operation
- Added `x_k8s_server_version` JMESPath function to fetch underlying cluster version
- Added `$config` built-in binding mapping to the underlying cluster client rest config

## 🔧 Fixes 🔧

- Fixed `x_k8s_resource_exists` JMESPath function not correctly processing `apiVersion`
- Fixed `x_k8s_exists` JMESPath function not correctly processing `apiVersion` and `kind`
- Fixed `x_k8s_list` JMESPath function not correctly exposing unstructured content
- Detect clustered resources in collectors to ignore namespace
- Fixed validation of checks in scripts
- Added missing validation of bindings

## 🎸 Misc 🎸

- Enabled cherry-pick bot
