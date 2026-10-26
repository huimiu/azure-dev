# Auth status

Auth status tells the user whether this `azd` config is signed in to Azure. A fresh scratch config is signed out. Signing in is a separate journey that needs a real Azure account.

## Sub-features

- `auth-help` lists `login`, `logout`, and `status`.
- `auth-status-text` reports that nobody is logged in.
- `auth-status-json` reports `status` `unauthenticated`.
- `auth-login` starts an Azure login. It is not part of the signed-out proof.

## How to get to it (user POV)

- Run `azd auth status` in a terminal.
- Run `azd auth status --output json` in a terminal.
- Run `azd auth --help` to see `login`, `status`, and `logout`.
- Run `azd auth login` to sign in. That opens a browser or device-code flow unless other login flags are supplied.

## Driving it with azd-verify

Preconditions:

- `azd-verify doctor` reports `status: healthy` for this `VERIFY_AZD_RUN_ID`.
- The scratch config has no saved account and `auth.useAzCliAuth` is unset.
- The harness has cleared `AZD_AUTH_ENDPOINT`, `AZD_AUTH_KEY`, and `AZD_AUTH_CERT`. Do not export them for this recipe.

- **Help.** Run `.cursor/skills/verify-azd/helpers/azd-verify drive --artifact auth/help -- auth --help`. Exit code `0`. `stdout.txt` contains `azd auth`, `login`, and `status`. `stderr.txt` is empty.
- **Text status.** Run `.cursor/skills/verify-azd/helpers/azd-verify drive --artifact auth/status-text -- auth status`. Exit code `0`. `stdout.txt` contains `Not logged in, run` and `azd auth login` and `to login to Azure`. `stderr.txt` is empty.
- **JSON status.** Run `.cursor/skills/verify-azd/helpers/azd-verify drive --artifact auth/status-json -- auth status --output json`. Exit code `0`. `stdout.txt` is JSON whose `status` is `unauthenticated`, with no `email` and no `clientId`. `stderr.txt` is empty.
- **Login.** Do not run `azd auth login` in this recipe. The unmet precondition is an Azure account and an interactive or device-code login against that account. A signed-out status is not a successful login.
- **Proof.** Keep the help and both status records. They prove the signed-out state only.

## Gotchas

- `azd auth login` without a real Azure account is not a passing journey. Do not weaken it to `auth status` and call login proven.
- Setting `auth.useAzCliAuth` makes `azd auth status` call the Azure CLI. This recipe expects built-in auth and the signed-out sentence above. A message that mentions `az login` means that setting was on.
- `AZD_AUTH_ENDPOINT` together with `AZD_AUTH_KEY` sends status to an external token endpoint. The harness strips both. If a drive is run outside the harness and those variables are set, the result is not this signed-out check.
- JSON status uses `unauthenticated`. The text status does not print that word; it prints the `Not logged in` sentence.
- `azd auth token` and `azd auth logout` need a signed-in account. They are not covered by the signed-out records.
