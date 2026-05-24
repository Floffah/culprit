#!/usr/bin/env bash
set -euo pipefail

## Recipe: Maculprit

# Reason: Remove config files and installed recipes when uninstalling maculprit.
# Matcher: $HOME/.maculprit
# Size: 37.0 B
rm -rf -- '/Users/ramsay/.maculprit'

