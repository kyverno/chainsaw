# Install

You can install the pre-compiled binary (in several ways), or compile from source.

## Install using `go install`

You can install with `go install` with:

```bash
go install github.com/kyverno/chainsaw@latest
```

## Manually

Download the pre-compiled binaries from the [releases page](https://github.com/kyverno/chainsaw/releases) and copy them to the desired location.

## Install using Docker

Chainsaw is also available as a Docker image which you can pull and run:

```bash
docker pull ghcr.io/kyverno/chainsaw:v0.0.2
```

> Note: Since Chainsaw relies on files for its operation (like test definitions), you will need to bind mount the necessary directories when running it via Docker.

```bash
docker run -d -v /path/on/host:/path/in/container ghcr.io/kyverno/chainsaw:v0.0.2 <chainsaw-command>
```

## Compile from sources

**clone:**

```bash
git clone github.com/kyverno/chainsaw
```

**get the dependencies:**

```bash
go mod tidy
```

**build:**

```bash
make build
```

**verify it works:**

```bash
./chainsaw version
```
