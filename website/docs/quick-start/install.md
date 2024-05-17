# Installation

You can install the pre-compiled binary (in several ways), compile from sources, or run with Docker.

We also provide a [GitHub action](#github-action) to easily install Chainsaw in your workflows.

## Install the pre-compiled binary

### Homebrew tap

**add tap:**

```bash
brew tap kyverno/chainsaw https://github.com/kyverno/chainsaw
```

**install chainsaw:**

```bash
brew install kyverno/chainsaw/chainsaw
```

!!! warning "Don't forget to specify the tap name"
    Homebrew core already has a tool named `chainsaw`.
    
    **Be sure that you specify the tap name when installing to install the right tool.**

### Manually

Download the pre-compiled binaries for your system from the [releases page](https://github.com/kyverno/chainsaw/releases) and copy them to the desired location.

### Install using `go install`

You can install with `go install` with:

```bash
go install github.com/kyverno/chainsaw@latest
```

## Run with Docker

Chainsaw is also available as a Docker image which you can pull and run:

```bash
docker pull ghcr.io/kyverno/chainsaw:<version>
```

!!! info

    Since Chainsaw relies on files for its operation (like test definitions), you will need to bind mount the necessary directories when running it via Docker.

```bash
docker run --rm                             \
    -v ./testdata/e2e/:/chainsaw/           \
    -v ${HOME}/.kube/:/etc/kubeconfig/      \
    -e KUBECONFIG=/etc/kubeconfig/config    \
    --network=host                          \
    ghcr.io/kyverno/chainsaw:<version>      \
    test /chainsaw --config /chainsaw/config.yaml
```

## Compile from sources

**clone:**

```bash
git clone https://github.com/kyverno/chainsaw.git
```

**build the binaries:**

```bash
cd chainsaw
go mod tidy
make build
```

**verify it works:**

```bash
./chainsaw version
```

## Install using Nix Package

To install `kyverno-chainsaw`, refer to the [documentation](https://search.nixos.org/packages?channel=unstable&show=kyverno-chainsaw&from=0&size=50&sort=relevance&type=packages&query=kyverno-chainsaw).

### On NixOS

```bash
nix-env -iA nixos.kyverno-chainsaw
```

### On Non-NixOS

```bash
nix-env -iA nixpkgs.kyverno-chainsaw
```

!!! warning
    Using nix-env permanently modifies a local profile of installed packages. This must be updated and maintained by the user in the same way as with a traditional package manager, foregoing many of the benefits that make Nix uniquely powerful. Using nix-shell or a NixOS configuration is recommended instead. 

### Using NixOS Configuration

Add the following Nix code to your NixOS Configuration, usually located in `/etc/nixos/configuration.nix` :

```nix
environment.systemPackages = [
  pkgs.kyverno-chainsaw
];
```

### Using nix-shell

A nix-shell will temporarily modify your `$PATH` environment variable. This can be used to try a piece of software before deciding to permanently install it. Use the following command to install `kyverno-chainsaw` :

```bash
nix-shell -p kyverno-chainsaw
```

## GitHub action

A GitHub action is available to install Chainsaw in your workflows.
See the [GitHub action](../gh-action.md) dedicated documentation.
