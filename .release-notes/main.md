# Release notes

Release notes for `TODO`.

## 💫 New features 💫

- Added bindings support when referencing step templates
- Restructured the reports system
- Added reports in JSON format back
- Added new JUnit report formats with different granularity (`JUNIT-TEST`, `JUNIT-STEP`, `JUNIT-OPERATION`)
- `XML` report format should be considered deprecated and replaced with one of the `JUNIT-*` format

## 🔧 Fixes 🔧

- Fixed default report name not set correctly
- Made `template` required when referencing a step template

## 📚 Docs 📚

- Started a JMESPath function example page in the docs (only `x_k8s_get` is available at this time)

## 🎸 Misc 🎸

- Removed test `v1alpha2` API for now
