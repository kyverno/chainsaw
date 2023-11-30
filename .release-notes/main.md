# Release notes

Release notes for `TODO`.

<!--

## :sparkles: UI changes :sparkles:

## :star: Examples :star:

## :boat: Tutorials :boat:

## :books: Docs :books:
-->

## :bangbang: Breaking chages :bangbang:

- Changed `check` to `expect` in delete `operation`
- Kuttl migration command split into `chainsaw kuttl migrate tests` and `chainsaw kuttl migrate config`

## :dizzy: New features :dizzy:

- Added support to register Chainsaw specific JMESPath functions for use in assertion trees
- Added inline manifest support to `assert` and `error` operations
- Added a `description` field to `Test`, `TestStep`, `Operation`, `Catch` and `Finally`
- Added a command to generate tests documentation

## :wrench: Fixes :wrench:

- Don't show `command` or `script` logs if there's nothing to show
- Fixed `timeouts` always set when migrating kuttl tests

## :guitar: Misc :guitar:

- We now have a [dedicated slack channel](https://kubernetes.slack.com/archives/C067LUFL43U)