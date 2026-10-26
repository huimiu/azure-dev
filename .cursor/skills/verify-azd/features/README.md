# azd verification map

This directory is the maintained source for verifying the user-facing behavior of the Azure Developer CLI (`azd`). Read the index before driving the CLI, then use the matching feature file as the recipe.

The primary surface is the `azd` CLI. The VS Code extension under `ext/vscode` is not a feature in this map.

## Baseline preconditions

- Build and launch with `.cursor/skills/verify-azd/helpers/azd-verify` so the binary is `/tmp/verify-azd/$VERIFY_AZD_RUN_ID/bin/azd`.
- Set `VERIFY_AZD_RUN_ID` to a new id per run. The harness sets `AZD_CONFIG_DIR` to that run's scratch `config/` directory.
- Run `azd-verify doctor` and require `status: healthy`, the scratch binary, the scratch config directory, and the stamped version.
- Never drive an `azd` that uses `~/.azd` or a binary installed by `go tool mage dev:install`.

## Driving conventions

- Start every recipe from a fresh launch unless its preconditions say otherwise.
- Treat every command as literal. Keep quoted names and flags unchanged.
- Run commands through `.cursor/skills/verify-azd/helpers/azd-verify drive -- <azd arguments>`.
- Pass `--cwd` only when the recipe names a project directory. Relative cwd values are created under the scratch `work/` directory.
- Restore nothing in the shared user home. Scratch state goes away in cleanup. Do not remove proof artifacts during cleanup.

## Proof and skip reporting

- Capture the user action and the resulting state, not only the final screen.
- CLI proof includes `command.txt`, `stdout.txt`, `stderr.txt`, `exit_code.txt`, and `transcript.txt` under `.cursor/skills/verify-azd/artifacts/<feature-id>/`.
- Mutation proof includes a `--snapshot` of the written file and a later read-only command that shows the same value.
- Record the feature id and entry point in the artifact path.
- Report an unreachable path with the attempted command and the unmet precondition.
- Do not report a skipped entry point as verified through a different path.

## Feature entry contract

Each feature file starts with an H1 title and one paragraph describing the user-visible behavior. It then uses exactly four H2 sections in this order.

1. `Sub-features` lists short IDs with one line for each behavior.
2. `How to get to it (user POV)` lists every user entry point.
3. `Driving it with azd-verify` starts with `Preconditions:` and uses labeled bullets that pair each user action with an exact command and observable result.
4. `Gotchas` lists traps that can waste or invalidate a verification run.

Keep implementation details out of the map. Name only user paths, stable handles, required state, commands, and observable proof.

## Features

- [Version and help](./version-help.md) covers `azd version`, JSON version output, root help, and version help.
- [Config](./config.md) covers showing, setting, reading, and unsetting `defaults.location` in an isolated config file.
- [Init](./init.md) covers `azd init --minimal` in an empty directory, including the project file and git repository it creates.
- [Template list](./template-list.md) covers `azd template list` and the network fetch it performs.
- [Auth status](./auth.md) covers signed-out `azd auth status`. Signing in needs a real Azure account.
