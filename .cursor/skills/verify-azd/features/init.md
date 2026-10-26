# Init

Init with `--minimal` creates a local `azure.yaml` project in an empty directory and initializes a git repository there. It does not provision Azure resources.

## Sub-features

- `init-minimal` accepts the directory name as the project name when `--no-prompt` is set.
- `init-files` writes `azure.yaml` and a `.gitignore` that ignores `.azure`.
- `init-git` creates a git repository in that directory.
- `init-help` shows the `--minimal` flag in `azd init --help`.

## How to get to it (user POV)

- From an empty directory, run `azd init --minimal`.
- From an empty directory, run `azd init --minimal --no-prompt` to accept the directory name.
- Run `azd init --help` to see `--minimal`, `--from-code`, and `--template`.

## Driving it with azd-verify

Preconditions:

- `azd-verify doctor` reports `status: healthy` and `git:` is a path, not `absent`.
- `git` is on `PATH`. The harness does not install git.
- The project directory `verify-init` does not already exist under the scratch `work/` directory.

- **Help.** Run `.cursor/skills/verify-azd/helpers/azd-verify drive --artifact init/help -- init --help`. Exit code `0`. `stdout.txt` contains `azd init` and `--minimal`. `stderr.txt` is empty.
- **Minimal init.** Run `.cursor/skills/verify-azd/helpers/azd-verify drive --cwd verify-init --artifact init/minimal --snapshot work/verify-init/azure.yaml --snapshot work/verify-init/.gitignore -- init --minimal --no-prompt`. Exit code `0`. `stdout.txt` contains `Generated azure.yaml project file.` and `Initialized git repository`. `stderr.txt` is empty.
- **Project file.** The `azure.yaml` snapshot contains `name: verify-init`.
- **Git ignore.** The `.gitignore` snapshot contains a line `.azure`.
- **Git repository.** After the drive, `git -C /tmp/verify-azd/$VERIFY_AZD_RUN_ID/work/verify-init rev-parse --is-inside-work-tree` prints `true`. Record that command and its stdout next to the drive record.
- **Proof.** The help record, the drive record, both snapshots, and the git rev-parse output together show the user command and the files left on disk.

## Gotchas

- This journey needs `git`. If doctor prints `git: absent`, stop and report that unmet precondition. Do not skip git initialization and still call init proven.
- `--minimal` and `--template` together fail. Use only `--minimal` here.
- A directory that already contains `azure.yaml` fails with a project-already-initialized error. Start from a new `--cwd`.
- `--no-prompt` uses the directory base name as the project name. The name `verify-init` is valid. A name with spaces or uppercase letters is a different journey.
- `azd init` without `--minimal`, `--from-code`, or `--template` prompts for an init mode. In a PTY that prompt waits for input. This recipe does not answer that prompt.
- `azd init --template` downloads a template and can create Azure environment files. It is not this feature and it is not proven by `--minimal`.
- Minimal init does not log in to Azure and does not prove `azd up`.
