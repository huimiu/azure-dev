# Template list

Template list shows sample application templates the user can pass to `azd init --template`. On a fresh config it registers the Awesome azd source and downloads that catalog.

## Sub-features

- `template-list-help` shows the list command and its filter flags.
- `template-list-json` prints the fetched catalog as JSON.
- `template-list-source` limits the list to the `awesome-azd` source.

## How to get to it (user POV)

- Run `azd template list` in a terminal.
- Run `azd template list --output json` for a stable list.
- Run `azd template list --source awesome-azd --output json` to keep the built-in catalog out of the result.
- Run `azd template list --help` to see `--source` and `--filter`.

## Driving it with azd-verify

Preconditions:

- `azd-verify doctor` reports `status: healthy` for this `VERIFY_AZD_RUN_ID`.
- This run has not already set config values that the proof must keep. The list command writes template sources into `config.json`.
- The machine can open HTTPS to `https://aka.ms/awesome-azd/templates.json`. No Azure subscription is required. No Azure login is required.

- **Help.** Run `.cursor/skills/verify-azd/helpers/azd-verify drive --artifact template-list/help -- template list --help`. Exit code `0`. `stdout.txt` contains `azd template list` and `--source`. `stderr.txt` is empty. This step does not fetch the catalog.
- **List.** Run `.cursor/skills/verify-azd/helpers/azd-verify drive --timeout 180 --artifact template-list/json --snapshot config/config.json -- template list --source awesome-azd --output json`. Exit code `0`. `stdout.txt` contains a JSON array of template objects, and at least one object has a repository path. `stderr.txt` may contain the `Retrieving templates` spinner. The snapshot shows an `awesome-azd` template source.
- **Proof.** Keep the help record, the JSON record, and the config snapshot. A non-zero exit or an error about the template source means the fetch failed. Record that output and stop. Do not replace it with an empty list or with help text.

## Gotchas

- The default list, without `--source`, still contacts Awesome azd and fails the whole command when that download fails. `--source awesome-azd` does not avoid the network.
- The first list on an empty config writes `template.sources` into `config.json` before the download finishes. A failed download can still leave that write behind. Do not run this recipe in the same run as the config feature.
- Table output is the default and includes gallery links. Use `--output json` for the proof.
- This command does not prove `azd init --template`. Cloning a template is a separate journey and needs network plus a destination directory.
- There is no offline catalog switch in this recipe. If HTTPS to `https://aka.ms/awesome-azd/templates.json` is blocked, the journey is unmet, not passed.
