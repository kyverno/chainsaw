# Kuttl Test Migration to Chainsaw

## Overview

The `chainsaw kuttl migrate` command is designed for the migration of KUTTL tests to Chainsaw. When executed, the command looks for KUTTL-defined tests and attempts to convert these into the equivalent Chainsaw-defined tests.

### Usage

```bash
chainsaw kuttl migrate [flags]
```

### Options

```bash
  -h, --help        help for migrate
      --overwrite   If set, overwrites original file.
      --save        If set, converted files will be saved.
```

## Description

On invocation, the command:

- Discovers folders from the specified paths.
- Reads the files within these folders, specifically looking for YAML files.
- Inspects each YAML file to check if it's a KUTTL resource. If it is, the command tries to convert it to a Chainsaw resource.
- The conversion handles two types of KUTTL resources: TestSuite and TestStep. It also reports an error for unsupported resources like `TestAssert`.
- If the `--save` flag is set, the converted Chainsaw tests are saved to a new file with the extension `.chainsaw.yaml`.

### Implementation Details

- **Discover Folders**: The command finds folders in the specified paths.
- **File Inspection**: It filters out non-YAML files and directories, focusing only on YAML files which might contain KUTTL test definitions.
- **Resource Conversion**: For identified KUTTL resources:
  - `TestSuite` is converted to Chainsaw's `Configuration` resource.
  - `TestStep` is converted to Chainsaw's `TestStep` resource.

  > **Note**: If the file contains a `TestAssert`, an error is reported since its conversion isn't currently supported.

- **Save Converted Tests**: If the `--save` flag is provided and if any resource within the YAML file needs saving (as determined by the migration process), the converted tests are saved. The file path for saving is determined by the `--overwrite` flag; if it is set, the original file will be overwritten, else a new `.chainsaw.yaml` file will be created

### SEE ALSO

- [chainsaw kuttl](chainsaw_kuttl.md)- Work with KUTTL tests
