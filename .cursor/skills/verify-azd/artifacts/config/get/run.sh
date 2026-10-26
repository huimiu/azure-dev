#!/usr/bin/env bash
set -u
cd /tmp/verify-azd/config-proof-0777/work || exit 97
fifo=/tmp/verify-azd/config-proof-0777/drives/03/stderr.fifo
stderr_file=/tmp/verify-azd/config-proof-0777/drives/03/stderr.txt
rm -f "${fifo}"
mkfifo "${fifo}"
tee "${stderr_file}" < "${fifo}" >&2 &
tee_pid=$!
ec=0
env -i \
  HOME=/home/ubuntu \
  PATH=/usr/local/cargo/bin:/exec-daemon:/exec-daemon/tools:/usr/local/cargo/bin:/home/ubuntu/.nvm/versions/node/v22.23.3/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin \
  TERM=xterm-256color \
  USER=ubuntu \
  LOGNAME=ubuntu \
  LANG=en_US.UTF-8 \
  LC_ALL=C.UTF-8 \
  TMPDIR=/tmp \
  COLUMNS=120 \
  LINES=40 \
  AZD_CONFIG_DIR=/tmp/verify-azd/config-proof-0777/config \
  AZURE_DEV_COLLECT_TELEMETRY=no \
  NO_COLOR=1 \
  AZD_SKIP_UPDATE_CHECK=true \
  AZD_DISABLE_AGENT_DETECT=1 \
  /tmp/verify-azd/config-proof-0777/bin/azd \
  config \
  get \
  defaults.location \
  2>"${fifo}" || ec=$?
wait "${tee_pid}" || true
printf '%s\n' "${ec}" > /tmp/verify-azd/config-proof-0777/drives/03/exit_code.txt
exit "${ec}"
