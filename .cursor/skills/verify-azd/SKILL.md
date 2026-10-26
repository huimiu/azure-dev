---
name: verify-azd
description: "Drive the Azure Developer CLI (azd) as a user would: build an isolated binary, health-check it, and run commands in a tmux PTY with a disposable AZD_CONFIG_DIR. Use to prove CLI behavior for version, help, config, init, template list, or auth status without touching the shared ~/.azd config."
---

# Verify azd

The primary user surface is the `azd` CLI. The VS Code extension under `ext/vscode` shells out to that CLI and is out of scope for this skill. Functional tests in `cli/azd/test/azdcli` and `cli/azd/test/functional` spawn the binary with pipes; this skill does not run `go test`. It builds a binary and drives it in its own tmux session, which is the same shape as a user typing at a terminal.

Each run has its own scratch directory `/tmp/verify-azd/$VERIFY_AZD_RUN_ID` containing the binary, config directory, and working directory. Two runs with different ids do not share state. Never point `AZD_CONFIG_DIR` at `~/.azd`, and never install the verification binary with `go tool mage dev:install` (that writes `azd-dev` into `~/.azd/bin`).

## Launch

Prerequisites on the machine: Go 1.26.1 or newer, `git`, `tmux`, and `python3`. From a clean run id, launch builds `azd` with the same version ldflags as `go tool mage dev:install` (`cli/version.txt` plus `-dev`, and the current commit) and writes the binary to the scratch `bin/azd`. It also creates empty `config/` and `work/` directories. Ready means the binary exists and `doctor` reports `status: healthy`. A successful `go build` alone is not ready.

```bash
export VERIFY_AZD_RUN_ID=config-proof-0777
.cursor/skills/verify-azd/helpers/azd-verify launch --artifact launch
```

Launch prints `status: launched`, the binary path, the config directory, and the version it stamped. Re-running launch for the same id fails until `cleanup`, unless you pass `--force` (that deletes the scratch dir for that id only).

The child environment for every later `azd` invocation is built with `env -i`. It keeps `HOME`, `PATH`, `TERM`, and locale, and sets:

- `AZD_CONFIG_DIR` to the scratch `config/` directory
- `AZURE_DEV_COLLECT_TELEMETRY=no`
- `NO_COLOR=1`
- `AZD_SKIP_UPDATE_CHECK=true`
- `AZD_DISABLE_AGENT_DETECT=1`

It does not inherit `AZD_AUTH_ENDPOINT`, `AZD_AUTH_KEY`, `AZD_AUTH_CERT`, `AZD_DEBUG`, `AZD_UI_PROMPT_ENDPOINT`, or a parent `AZD_CONFIG_DIR`. `AZD_DISABLE_AGENT_DETECT=1` keeps agent markers such as `CURSOR_AGENT` from turning on `--no-prompt` by themselves. Pass `--no-prompt` only when the feature recipe says to.

## Doctor

Doctor is read-only with respect to the run config directory and `~/.azd`. It does not start a tmux session. It checks that the binary is the scratch binary, that the config directory is the scratch `config/` and not `~/.azd`, that `tmux` and `python3` are available, and that `azd version` and `azd version --output json` match the identity launch stamped. Those two version commands run with `AZD_CONFIG_DIR` pointed at a throwaway directory under the scratch tree, which doctor deletes before it returns. `azd version --output json` writes the default extension registry into whatever config directory it is given, so doctor must not point that command at the run config. Text output must be `azd version <version> (commit <40-hex>) (stable)` unless the stamped version contains `daily.`, in which case the suffix is ` (daily)`. JSON `azd.version` is the semver only and `azd.commit` is the same commit, with empty stderr. If `~/.azd` already exists, doctor records its directory mtime and fails when that mtime changes. It prints `git:` as the path to `git` or `absent`; a missing `git` does not fail doctor, and it does block `azd init`.

```bash
.cursor/skills/verify-azd/helpers/azd-verify doctor --artifact doctor
```

Drive only when the report contains `status: healthy`. Exit code `1` means unhealthy. The `reason:` lines are the failure.

## Drive

`drive` starts a tmux session named `verify-azd-$VERIFY_AZD_RUN_ID-dNN` (120 columns by 40 rows). The session waits until `remain-on-exit` is set, then runs `azd` with stdin and stdout attached to the PTY. Stderr is tee'd to `stderr.txt` and back onto the terminal. Arguments after `--` are `azd` arguments. A leading `azd` token is stripped so the recipe can show the user command.

```bash
.cursor/skills/verify-azd/helpers/azd-verify drive \
  --artifact config/set \
  --snapshot config/config.json \
  -- config set defaults.location eastus2
```

- `--cwd` is a directory under the scratch tree. A relative path is created under `work/`. Refusing a path outside the scratch tree is intentional.
- `--timeout` is seconds (default 90). On timeout the exit code is `124` and the pane is still captured.
- `--snapshot` copies a file from the scratch tree into the evidence directory. The path is relative to the scratch root, for example `config/config.json` or `work/verify-init/azure.yaml`. A missing snapshot fails the drive.
- `--artifact` copies the drive record under `.cursor/skills/verify-azd/artifacts/<relative-path>`.

The drive record contains `command.txt`, `stdout.txt` (pane text with the tmux `Pane is dead` status line and trailing blank screen rows removed), `stderr.txt`, `transcript.txt` (stdout section then stderr section), `exit_code.txt`, and `meta.txt`. `drive` leaves the tmux session in place for `cleanup`. It exits with `azd`'s exit code.

Interactive prompts only work when the recipe does not pass `--no-prompt`, because stdin is the PTY. Do not drive a second run's session, and do not type into a shared user terminal.

## Evidence

Proof is the user command plus the resulting state. For every drive, keep `command.txt`, `stdout.txt`, `stderr.txt`, `exit_code.txt`, and `transcript.txt`. When the command writes a file, snapshot that file in the same drive record. Doctor evidence is `doctor.txt` plus the version stdout, stderr, and JSON captures. Launch evidence is `launch.txt` and `build.log`.

Write evidence only under `.cursor/skills/verify-azd/artifacts/`. A proof of one feature does not cover the other feature files. Record the feature id in the artifact path (`config/set`, `version-help/text`, and so on). Do not call a journey proven because a different command succeeded.

`azd` has no dry-run switch that this skill uses. `AZURE_DEV_COLLECT_TELEMETRY=no` is the supported opt-out and skips telemetry initialization. `AZD_SKIP_UPDATE_CHECK=true` skips the network update check. Confirm those by empty stderr on `azd version` and by the doctor mtime check on `~/.azd`, not by the variable name alone.

Commands that need Azure are listed in the feature gotchas. Do not record them as proven unless the evidence shows a real Azure response from this run.

## Cleanup

Cleanup kills tmux sessions named for this run id and deletes `/tmp/verify-azd/$VERIFY_AZD_RUN_ID`. It does not kill other sessions, it does not delete `~/.azd`, and it does not delete `.cursor/skills/verify-azd/artifacts`.

```bash
.cursor/skills/verify-azd/helpers/azd-verify cleanup
```

After cleanup, the evidence files must still be at the artifact paths. A missing proof file means cleanup or the drive target was wrong; fix that before calling the run done.

Run cleanup after a failed launch or drive too, so a broken attempt does not leave sessions behind. Use the same `VERIFY_AZD_RUN_ID`.

## Helpers

The only helper is `.cursor/skills/verify-azd/helpers/azd-verify`. It is executable. Subcommands are `launch`, `doctor`, `drive`, and `cleanup`. Invoke it from the repository root as shown above. `azd-verify help` prints the same usage.

## Maintenance

When commands, flags, or required environment change, update this skill and the feature map together. Use `/maintain-verification-skill` to keep the map aligned with the CLI. Do not mark a feature proven from a unit test or from help text alone.
