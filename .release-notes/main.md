# Release notes

Release notes for `TODO`.

<!--
## ‼️ Breaking changes ‼️

## ✨ UI changes ✨

## ⭐ Examples ⭐

## ⛵ Tutorials ⛵
-->

## 💫 New features 💫

- Added support for running chainsaw without a Kubernetes cluster
- Automatically add json schemas when converting tests from kuttl to chainsaw
- Support glob patterns in operation `file` i.e. filepath
- Added support for passing arbitrary values when invoking `chainsaw test`

## 🔧 Fixes 🔧

- Fixed a couple of kuttl to chainsaw migration bugs

## 📚 Docs 📚

- Added community meetings docs
- Added google groups in community docs and README
- Removed all references to `TestStep` in the docs as this is not supported anymore
- Remove TestStep reference examples

## 🎸 Misc 🎸

- Added kubebuilder community recording link where chainsaw was presented
- Added a `MAINTAINERS.md` file
- Added makefile target to check go code coverage
- Added CI to check PR semantics
