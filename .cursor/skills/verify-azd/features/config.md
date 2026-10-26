# Config

Config lets a user store default prompt values such as Azure location in `config.json` and read them back with a later command.

## Sub-features

- `config-show-empty` prints an empty configuration before any value is set.
- `config-set` stores `defaults.location`.
- `config-get` prints that stored location.
- `config-show-value` prints the whole file as JSON.
- `config-file` is the `config.json` written under the scratch config directory.
- `config-unset` removes the location and leaves an empty configuration.

## How to get to it (user POV)

- Run `azd config show` in a terminal.
- Run `azd config set defaults.location <location>` in a terminal.
- Run `azd config get defaults.location` in a terminal.
- Run `azd config unset defaults.location` in a terminal.
- Run `azd config --help` to see the `set`, `get`, `show`, and `unset` commands.

## Driving it with azd-verify

Preconditions:

- `azd-verify doctor` reports `status: healthy` for this `VERIFY_AZD_RUN_ID`.
- The scratch config directory does not yet contain `config.json`.
- No earlier command in this run has written config. `azd template list` writes template sources into the same file.

- **Show empty.** Run `.cursor/skills/verify-azd/helpers/azd-verify drive --artifact config/show-empty -- config show`. Exit code `0`. `stdout.txt` is `{}` plus a newline. `stderr.txt` is empty. The scratch file `config/config.json` does not exist yet.
- **Set location.** Run `.cursor/skills/verify-azd/helpers/azd-verify drive --artifact config/set --snapshot config/config.json -- config set defaults.location eastus2`. Exit code `0`. `stdout.txt` and `stderr.txt` are empty. The snapshot `config.json` contains `"location": "eastus2"`.
- **Get location.** Run `.cursor/skills/verify-azd/helpers/azd-verify drive --artifact config/get -- config get defaults.location`. Exit code `0`. `stdout.txt` is `"eastus2"` plus a newline. `stderr.txt` is empty.
- **Show value.** Run `.cursor/skills/verify-azd/helpers/azd-verify drive --artifact config/show --snapshot config/config.json -- config show`. Exit code `0`. `stdout.txt` is JSON with `defaults.location` equal to `eastus2`. The snapshot matches that value.
- **Unset location.** Run `.cursor/skills/verify-azd/helpers/azd-verify drive --artifact config/unset -- config unset defaults.location`. Exit code `0`. `stdout.txt` and `stderr.txt` are empty.
- **Show after unset.** Run `.cursor/skills/verify-azd/helpers/azd-verify drive --artifact config/show-after-unset --snapshot config/config.json -- config show`. Exit code `0`. `stdout.txt` is `{}` plus a newline. The snapshot no longer contains `eastus2`.
- **Proof.** The `config/set` snapshot is the write, and `config/get` plus `config/show` are the read-back. Keep all six records.

## Gotchas

- `azd config set` prints nothing on success. An empty transcript with exit code `0` is not proof until `config.json` and `azd config get` show the value.
- `azd config get` and `azd config show` default to JSON. The location is a JSON string, including the quotes.
- The first config command creates the scratch config directory. It does not create `config.json` until `set`.
- `azd config reset` asks for confirmation unless `--force` is passed. This recipe uses `unset` so it does not prompt.
- `azd template list` on a fresh config writes `template.sources` before it fetches templates. Do not run it in the same run as this recipe.
- A value stored here is not an Azure subscription check. `eastus2` is only a local string.
