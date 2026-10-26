#!/usr/bin/env bash
set -eu
go_file="$1"
run_sh="$2"
while [[ ! -f "${go_file}" ]]; do
  sleep 0.05
done
exec bash "${run_sh}"
