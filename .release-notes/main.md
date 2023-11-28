# Release notes

Release notes for `TODO`.

<!--
## :sparkles: UI changes :sparkles:

## :star: Examples :star:

## :boat: Tutorials :boat:

## :wrench: Fixes :wrench:

## :guitar: Misc :guitar:
-->

## :bangbang: Breaking chages :bangbang:

- Object reference in `delete` is now under the `ref` field
- Check in `apply` and `create` operations was renamed to `expect` and is now an array of `Expectation`s (a combination of a match and a check)
- Additional data passed to `check`s are now done using bindings (`$error`, `$stdout`, `$stderr`, etc...)

## :dizzy: New features :dizzy:

- Added timeout support in `try` and `catch` handlers
- Added assertion tree check in `delete` operation
- Added a new configuration option to force termination graceful period on `Pod`, `Deployment`, `StatefulSet`, `DaemonSet`, `Job` and `CronJob`
- Added reports support
- Completed kuttl migration command with `TestAssert` support

## :wrench: Fixes :wrench:

- Fixed a kuttl migration failure in case of unsupported file name
- Fixed a potential invalid name when migrating a kuttl test step
- Fixed `check` set to `null` in kuttl migration command
- Fixed a manifest discovery issue where manifests could be loaded in the wrong order
- Fixed a manifest discovery issue where error manifests where not discovered correctly
- Fixed controller-runtime logger not initialised
- Fixed `command` and `script` bug not receiving the right error in `$error` check

## :books: Docs :books:

- Added missing descriptions on API reference docs

## :guitar: Misc :guitar:

- Added command line log in `command` and `script` operations
