# Release docs

This doc contains information for releasing a new version.

## Create a release

Creating a release can be done by pushing a tag to the GitHub repository (beginning with `v`).

The [release workflow](../../.github/workflows/release.yaml) will take care of creating the GitHub release and will publish artifacts.

```shell
VERSION="v0.1.0"
TAG=$VERSION

git tag $TAG -m "tag $TAG" -a
git push origin $TAG
```

## Release notes

Release notes for the `main` branch live in [main.md](../../.release-notes/main.md).

Make sure it is up to date and rename the file to the version being released.

You can then copy [_template.md](../../.release-notes/_template.md) to [main.md](../../.release-notes/main.md) for the next release.

## Publish documentation

Publishing the documentation for a release is decoupled from cutting a release.

To publish the documentation push a tag to the GitHub repository (beginning with `docs-v`).

```shell
VERSION="v0.1.0"
TAG=docs-$VERSION

git tag $TAG -m "tag $TAG" -a
git push origin $TAG
```

## Publish GitHub action

Once the release is cut, bump the default version in the GH action [here](https://github.com/kyverno/action-install-chainsaw/blob/main/action.yml).

Publish a new version of the GitHub action.