# Release docs

This docs contains informations for releasing releasing.

## Create a release

Creating a release can be done by pushing a tag to the GitHub repository (begining with `v`).

The [release workflow](../../.github/workflows/release.yaml) will take care of creating the GitHub release and will publish artifacts.

```shell
VERSION="v0.1.0"
TAG=$VERSION

git tag $TAG -m "tag $TAG" -a
git push origin $TAG
```

## Release notes

Release notes for the `main` branch lives in [main.md](../../.release-notes/main.md).

Make sure it is up to date and rename the file to the version being released.

You can then copy [_template.md](../../.release-notes/_template.md) to [main.md](../../.release-notes/main.md) for the next release.
