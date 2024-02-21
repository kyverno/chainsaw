# Testing command output

Chainsaw can be used to easily check terminal output from CLIs and other commands. This is useful in that convoluted bash scripts involving chaining together tools like `grep` can be avoided or at least minimized to only complex use cases. Output to both stdout and stderr can be checked for a given string or precise contents.

## Checking Output Contains

One basic use case for content checking is that the output simply contains a given string or piece of content. For example, you might want to run automated tests on a CLI binary you build to ensure that a given command produces output that contains some content you specify somewhere in the output. Let's use the following output from the `kubectl version` command to show these examples.

```sh
kubectl version

Client Version: v1.28.2
Kustomize Version: v5.0.4-0.20230601165947-6ce0bf390ce3
Server Version: v1.27.4+k3s1
```

Below is an example that ensures the string '1.28' is found somewhere in that output. So long as the content is present anywhere, the test will succeed. To perform this check, the [`contains()`](../jp/functions.md) JMESPath filter is used.

```yaml
apiVersion: chainsaw.kyverno.io/v1alpha1
kind: Test
metadata:
  name: test
spec:
  steps:
  - name: Check kubectl
    try:
    - script:
        content: kubectl version
        check:
          # This check ensures that the string '1.28' is found
          # in stdout or else fails
          (contains($stdout, '1.28')): true
```

Checks for content containing a given value can be negated as well. For example, checking to ensure the output does NOT contain the string '1.25'.

```yaml
- script:
    content: kubectl version
    check:
      # This check ensures that the string '1.25' is NOT found
      # in stdout or else fails
      (contains($stdout, '1.25')): false
```

## Checking Output Is Exactly

In addition to checking that CLI/command output contains some contents, you may need to ensure that the contents are exactly as intended. The Chainsaw test below accomplishes this by comparing the entire contents of stdout with those specified in the block scalar. If so much as one character, space, or line break is off, the test will fail. This is useful in that not only can content be checked but the formatting of that content can be ensured it matches a given declaration.

```yaml
apiVersion: chainsaw.kyverno.io/v1alpha1
kind: Test
metadata:
  name: test
spec:
  steps:
  - name: Check kubectl
    try:
    - script:
        content: kubectl version
        check:
          # This check ensures the contents of stdout are exactly as shown.
          # Any deviations will cause a failure.
          ($stdout): |-
            Client Version: v1.28.2
            Kustomize Version: v5.0.4-0.20230601165947-6ce0bf390ce3
            Server Version: v1.27.4+k3s1
```

## Checking Output In Errors

In addition to testing that commands succeed and with output in a given shape, it's equally valuable and necessary to perform negative tests; that tests fail and with contents that are as expected. Similarly, those checks can be for output which has some contents as well as output which appears exactly as desired. For example, you may wish to check that running the `kubectl foo` command not only fails as expected but that the output shown to users contains a certain word or sentence.

```sh
kubectl foo

error: unknown command "foo" for "kubectl"

Did you mean this?
        top
```

Below you can see an example where the command `kubectl foo` is expected to fail but that the error message returned contains some output, in this case the string 'top'.

```yaml
apiVersion: chainsaw.kyverno.io/v1alpha1
kind: Test
metadata:
  name: test
spec:
  steps:
  - name: Check bad kubectl command
    try:
    - script:
        content: kubectl foo
        check:
          # This checks that the result of the content was an error.
          ($error != null): true
          # This check below ensures that the string 'top' is found in stderr or else fails
          (contains($stderr, 'top')): true
```

Likewise, this failure output can be checked that it is precise. Note that in the example below, due to the use of a tab character in the output of `kubectl foo`, the value of the `($stderr)` field is given as a string to preserve these non-printing characters.

```yaml
apiVersion: chainsaw.kyverno.io/v1alpha1
kind: Test
metadata:
  name: test
spec:
  steps:
  - name: Check kubectl
    try:
    - script:
        content: kubectl foo
        check:
          # This checks that the result of the content was an error.
          ($error != null): true
          # This checks that the output is exactly as intended.
          ($stderr): "error: unknown command \"foo\" for \"kubectl\"\n\nDid you mean this?\n\ttop"
```
