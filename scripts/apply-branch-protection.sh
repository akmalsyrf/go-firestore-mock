#!/usr/bin/env bash
# apply-branch-protection.sh — solo-maintainer friendly protection for main.
# Requires: gh auth login (valid token with repo admin).
set -euo pipefail

REPO="${1:-akmalsyrf/go-firestore-mock}"
BRANCH="${2:-main}"

gh api "repos/${REPO}/branches/${BRANCH}/protection" --method PUT --input - <<EOF
{
  "required_status_checks": {
    "strict": true,
    "contexts": ["gate"]
  },
  "enforce_admins": false,
  "required_pull_request_reviews": null,
  "restrictions": null,
  "allow_force_pushes": false,
  "allow_deletions": false
}
EOF

echo "branch protection applied: ${REPO}@${BRANCH} requires status check 'gate', no required reviews"
