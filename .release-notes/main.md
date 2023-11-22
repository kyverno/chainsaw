# Release notes

Release notes for `TODO`.

## :bangbang: Breaking chages :bangbang:

- Replaced `shouldFail` field in `create` and `apply` operations by assertion trees to validate outputs and error

## :dizzy: New features :dizzy:

- Added assertion tree support in `command` and `script` operations to validate command outputs (stdout and stderr) and error
- Added support for `dryRun` in `create` and `apply` operations
- Added `--test-file` flag to use a custom test file name
- Added raw resource support in `Test` and `TestStep`s to provide a resource directly instead of referencing a file
- Improved loggs, added a `STATUS` column and operation details are now shown in sections

## :wrench: Fixes :wrench:

- Sort error messages to get predictable output
- Current working directory not set in `command` and `script` operations
- Fixed a few assertion tree bugs in [kyverno-json](https://github.com/kyverno/kyverno-json)
- Fixed an invalid regex check in GitHub action
- Fixed cases where cleanup timeout was not used correctly

## :books: Docs :books:

- Added multi version docs support

## :guitar: Misc :guitar:

- Added `command` and `script` operations support to KUTTL migration command
