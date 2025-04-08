# Shell Completion

Chainsaw provides shell completion support for Bash, Zsh, Fish, and PowerShell. Once set up, you can use the <Tab> key to auto-complete chainsaw commands, flags, and even some arguments, which significantly improves the command-line experience.

## Generating Completion Scripts

You can generate shell completion scripts using the `chainsaw completion` command:

```bash
# For Bash
chainsaw completion bash

# For Zsh
chainsaw completion zsh

# For Fish
chainsaw completion fish

# For PowerShell
chainsaw completion powershell
```

## Setting Up Completion

### Bash

To enable completion in your current Bash session:

```bash
source <(chainsaw completion bash)
```

To enable completion for all sessions, add the above line to your `~/.bashrc` file:

```bash
echo 'source <(chainsaw completion bash)' >> ~/.bashrc
```

Alternatively, you can save the completion script to the bash-completion directory:

```bash
# On Linux
chainsaw completion bash > /etc/bash_completion.d/chainsaw

# On macOS with Homebrew
chainsaw completion bash > $(brew --prefix)/etc/bash_completion.d/chainsaw
```

### Zsh

To enable completion in your current Zsh session:

```bash
source <(chainsaw completion zsh)
```

To enable completion for all sessions, add the above line to your `~/.zshrc` file:

```bash
echo 'source <(chainsaw completion zsh)' >> ~/.zshrc
```

Alternatively, you can save the completion script to a directory in your `$fpath`:

```bash
# Create a directory for completions if it doesn't exist
mkdir -p ~/.zsh/completion
# Generate and save the completion script
chainsaw completion zsh > ~/.zsh/completion/_chainsaw

# Make sure the directory is in your fpath by adding to ~/.zshrc:
echo 'fpath=(~/.zsh/completion $fpath)' >> ~/.zshrc
echo 'autoload -U compinit; compinit' >> ~/.zshrc
```

### Fish

To enable completion in Fish:

```bash
chainsaw completion fish > ~/.config/fish/completions/chainsaw.fish
```

### PowerShell

To enable completion in PowerShell:

```powershell
chainsaw completion powershell | Out-String | Invoke-Expression
```

To make it persistent, add the above line to your PowerShell profile:

```powershell
# Find the profile path
echo $PROFILE

# Add the completion command to your profile
chainsaw completion powershell | Out-String | Out-File -Append $PROFILE
```

## Testing Completion

After setting up completion, you can test it by typing `chainsaw` followed by a space and pressing <Tab>. This should show available subcommands. You can also try typing partial commands like `chainsaw te` and then pressing <Tab>, which should complete to `chainsaw test`.

## Detailed Reference

For more detailed information about each completion command, see the reference documentation:

- [Bash Completion](../reference/commands/chainsaw_completion_bash.md)
- [Fish Completion](../reference/commands/chainsaw_completion_fish.md)
- [PowerShell Completion](../reference/commands/chainsaw_completion_powershell.md)
- [Zsh Completion](../reference/commands/chainsaw_completion_zsh.md)
