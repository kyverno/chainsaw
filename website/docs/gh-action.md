# GitHub action

A GitHub action is available to easily install Chainsaw in your workflows.

The GitHub action is available at [kyverno/action-install-chainsaw](https://github.com/kyverno/action-install-chainsaw) or in the [marketplace](https://github.com/marketplace/actions/chainsaw-installer).

## Usage

This action currently supports GitHub-provided Linux, macOS and Windows runners (self-hosted runners may not work).

Add the following entry to your Github workflow YAML file:

```yaml
uses: kyverno/action-install-chainsaw@v0.1.0
with:
  release: v0.1.0 # optional
```

Example using a pinned version:

```yaml
jobs:
  example:
    runs-on: ubuntu-latest

    permissions: {}

    name: Install Chainsaw
    steps:
      - name: Install Chainsaw
        uses: kyverno/action-install-chainsaw@v0.1.0
        with:
          release: v0.0.9
      - name: Check install
        run: chainsaw version
```

Example using the default version:

```yaml
jobs:
  example:
    runs-on: ubuntu-latest

    permissions: {}

    name: Install Chainsaw
    steps:
      - name: Install Chainsaw
        uses: kyverno/action-install-chainsaw@v0.1.0
      - name: Check install
        run: chainsaw version
```

Example using [cosign](https://github.com/sigstore/cosign) verification:

```yaml
jobs:
  example:
    runs-on: ubuntu-latest

    permissions: {}

    name: Install Chainsaw
    steps:
      - name: Install Cosign
        uses: sigstore/cosign-installer@v3.1.1
      - name: Install Chainsaw
        uses: kyverno/action-install-chainsaw@v0.1.0
        with:
          verify: true
      - name: Check install
        run: chainsaw version
```

If you want to install Chainsaw from its main version by using `go install` under the hood, you can set `release` as `main`.
Once you did that, Chainsaw will be installed via `go install` which means that please ensure that go is installed.

Example of installing Chainsaw via `go install`:

```yaml
jobs:
  example:
    runs-on: ubuntu-latest

    permissions: {}

    name: Install Chainsaw via go install
    steps:
      - name: Install go
        uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      - name: Install Chainsaw
        uses: kyverno/action-install-chainsaw@v0.1.0
        with:
          release: main
      - name: Check install
        run: chainsaw version
```

### Optional Inputs

The following optional inputs:

| Input | Description |
| --- | --- |
| `release` | `chainsaw` version to use instead of the default. |
| `install-dir` | directory to place the `chainsaw` binary into instead of the default (`$HOME/.chainsaw`). |
| `use-sudo` | set to `true` if `install-dir` location requires sudo privs. Defaults to false. |
| `verify` | set to `true` to enable [cosign](https://github.com/sigstore/cosign) verification of the downloaded archive. |
