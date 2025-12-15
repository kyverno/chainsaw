# Sharding

If for some reason you need to distribute tests execution across different shards, Chainsaw can be instructed to run only a chunk of the discovered tests with 
`--shard-index` and `--shard-count` flags.

`--shard-count` defines the total number of shards and `--shard-index` defines the current shard.
Those flags are used to determine the chunk of tests that need to be executed.

## Configuration

!!! note
    Sharding can't be configured with a configuration file.

### With flags

```bash
chainsaw test       \
  --shard-count 4   \
  --shard-index 2
```
