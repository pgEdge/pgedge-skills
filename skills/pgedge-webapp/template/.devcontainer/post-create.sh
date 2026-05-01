#!/usr/bin/env bash
set -euo pipefail

# Install Helm (devcontainer features doesn't ship a Helm feature).
curl -fsSL https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 | bash

# Pre-install dependencies so the dev loop starts hot.
( cd server/src && go mod download )
( cd client && npm ci )
