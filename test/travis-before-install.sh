#!/bin/bash
set -o xtrace

# Travis does shallow clones, so there is no master branch present.
# But test-no-outdated-migrations.sh needs to check diffs against master.
# Fetch just the master branch from origin.
( git fetch origin master
git branch master FETCH_HEAD ) &
# Github-PR-Status secret
if [ -n "$encrypted_53b2630f0fb4_key" ]; then
  openssl aes-256-cbc \
    -K $encrypted_53b2630f0fb4_key -iv $encrypted_53b2630f0fb4_iv \
    -in test/github-secret.json.enc -out /tmp/github-secret.json -d
fi

./test/setup.sh

set +o xtrace
