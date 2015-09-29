#!/bin/bash
set -o xtrace

# Travis does shallow clones, so there is no master branch present.
# But test-no-outdated-migrations.sh needs to check diffs against master.
# Fetch just the master branch from origin.
time git fetch origin master
git branch master FETCH_HEAD
# Github-PR-Status secret
if [ -n "$encrypted_53b2630f0fb4_key" ]; then
  openssl aes-256-cbc \
    -K $encrypted_53b2630f0fb4_key -iv $encrypted_53b2630f0fb4_iv \
    -in test/github-secret.json.enc -out /tmp/github-secret.json -d
fi

time travis_retry go get   golang.org/x/tools/cmd/vet 
time travis_retry go get   golang.org/x/tools/cmd/cover 
time travis_retry go get   github.com/golang/lint/golint 
time travis_retry go get   github.com/mattn/goveralls 
time travis_retry go get   github.com/modocache/gover 
time travis_retry go get   github.com/jcjones/github-pr-status 
time travis_retry go get   github.com/letsencrypt/goose/cmd/goose

# Boulder consists of multiple Go packages, which
# refer to each other by their absolute GitHub path,
# e.g. github.com/letsencrypt/boulder/analysis. That means, by default, if
# someone forks the repo, Travis won't pass on their own repo. To fix that,
# we add a symlink.
mkdir -p $TRAVIS_BUILD_DIR $GOPATH/src/github.com/letsencrypt
if [ ! -d $GOPATH/src/github.com/letsencrypt/boulder ] ; then
  ln -s $TRAVIS_BUILD_DIR $GOPATH/src/github.com/letsencrypt/boulder
fi

set +o xtrace
