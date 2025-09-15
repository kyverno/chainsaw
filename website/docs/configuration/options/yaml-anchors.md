# YAML anchors

YAML anchors and pointers are not well defined across YAML version (https://ktomk.github.io/writing/yaml-anchor-alias-and-merge-key.html).
For example _Merge Keys_ are only part of YAML 1.1 which is deprecated, and not part of YAML 1.2

While the feature seems esoteric enough feature and should probably not be relied on, it is sometimes convenient.

Chainsaw supports the `--remarshal` flag to experimentally enable YAML Anchor, Aliases and Merge Keys before parsing the the YAML content.

## Configuration

!!! note
    Remarshaling can't be configured with a configuration file.

### With flags

```bash
chainsaw test --remarshal
```
