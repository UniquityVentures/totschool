#!/bin/sh
set -eu

root=$(CDPATH= cd -- "$(dirname "$0")/.." && pwd)
cd "$root"

git config core.hooksPath .githooks
chmod +x .githooks/pre-commit
echo "Git hooks enabled from $root/.githooks"
