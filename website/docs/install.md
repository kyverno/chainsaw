# Install

You can install the pre-compiled binary (in several ways), compile from sources, or run with Docker.

We also provide a [GitHub action](#github-action) to easily install Chainsaw in your workflows.

## Install the pre-compiled binary

### Homebrew tap

**add tap:**

```bash
$ brew tap kyverno/chainsaw https://github.com/kyverno/chainsaw
```

**install chainsaw:**

```bash
$ brew install kyverno/chainsaw/chainsaw
```

!!! warning "Don't forget to specify the tap name"
    Homebrew core already has a tool named `chainsaw`.
    
    **Be sure that you specify the tap name when installing to install the right tool.**

### Manually

Download the pre-compiled binaries for your system from the [releases page](https://github.com/kyverno/chainsaw/releases) and copy them to the desired location.

## Install using `go install`

You can install with `go install` with:

```bash
$ go install github.com/kyverno/chainsaw@latest
```

## Running with Docker

Chainsaw is also available as a Docker image which you can pull and run:

```bash
$ docker pull ghcr.io/kyverno/chainsaw:<version>
```

!!! info

    Since Chainsaw relies on files for its operation (like test definitions), you will need to bind mount the necessary directories when running it via Docker.

```bash
$ docker run --rm                       \
    -v /path/on/host:/path/in/container \
    ghcr.io/kyverno/chainsaw:<version>  \
    <chainsaw-command>
```

## Compiling from sources

**clone:**

```bash
$ git clone https://github.com/kyverno/chainsaw.git
```
**build the binaries:**

```bash
$ cd chainsaw
$ go mod tidy
$ make build
```

**verify it works:**

```bash
$ ./chainsaw version
```

## GitHub action

A GitHub action is available to install Chainsaw in your workflows.
See the [GitHub action](./gh-action.md) dedicated documentation.