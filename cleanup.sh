#!/usr/bin/env bash
set -euo pipefail

## Recipe: Culprit

# Reason: Remove config files and installed recipes when uninstalling culprit.
# Matcher: $HOME/.culprit
# Size: 37.0 B
rm -rf -- '/Users/ramsay/.culprit'
