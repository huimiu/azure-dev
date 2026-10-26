# Version and help

Version and help tell the user which `azd` binary they are running and which commands that binary exposes, without changing config or project files.

## Sub-features

- `version-text` prints the stamped version and channel on stdout.
- `version-json` prints the same version and commit as JSON.
- `help-root` lists the built-in command groups from `azd --help`.
- `help-version` shows the usage line for `azd version`.

## How to get to it (user POV)

- Run `azd version` in a terminal.
- Run `azd version --output json` in a terminal.
- Run `azd --help` or `azd help` in a terminal.
- Run `azd version --help` in a terminal.

## Driving it with azd-verify

Preconditions:

- `azd-verify doctor` reports `status: healthy` for this `VERIFY_AZD_RUN_ID`.
- The scratch config directory has no `config.json` yet.

- **Text version.** Run `.cursor/skills/verify-azd/helpers/azd-verify drive --artifact version-help/text -- version`. Exit code `0`. `stdout.txt` is `azd version <version> (stable)` when the stamped version does not contain `daily.`. `stderr.txt` is empty. `<version>` is the `version:` line from doctor.
- **JSON version.** Run `.cursor/skills/verify-azd/helpers/azd-verify drive --artifact version-help/json -- version --output json`. Exit code `0`. `stdout.txt` is JSON whose `azd.version` equals the doctor version and whose `azd.commit` equals the doctor commit. `stderr.txt` is empty.
- **Root help.** Run `.cursor/skills/verify-azd/helpers/azd-verify drive --artifact version-help/root-help -- --help`. Exit code `0`. `stdout.txt` contains `Azure Developer CLI`, the `init` command, the `config` command, the `auth` command, the `template` command, and the `version` command. `stderr.txt` is empty.
- **Version help.** Run `.cursor/skills/verify-azd/helpers/azd-verify drive --artifact version-help/version-help -- version --help`. Exit code `0`. `stdout.txt` contains `azd version [flags]` and `Print the version number of Azure Developer CLI.` `stderr.txt` is empty.
- **Proof.** Keep the four drive records. None of them create `config/config.json` in the scratch tree.

## Gotchas

- `azd version --docs` opens a browser. Do not use it as the help check.
- Root help includes an extensions block only for extensions installed in that config directory. A fresh scratch config has no extensions; do not require extension names in `stdout.txt`.
- `go build` without the launch ldflags prints `azd version 0.0.0-dev.0 (stable)`. Doctor rejects that binary because launch stamps `cli/version.txt` plus `-dev`.
- Help text is not a substitute for running `version`. Capture both the text and JSON commands.
