#!/bin/bash

SCRIPTDIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

# Ensure gopath is set and go is installed
if [[ ! -d $GOPATH ]] || [[ ! -d $GOBIN ]] || [[ ! -x "$(which go)" ]]; then
  echo "Your \$GOPATH is not set or go is not installed,"
  echo "ensure you have a working installation of go before trying again..."
  echo "https://golang.org/doc/install"
  exit 1
fi

SDK_REPO="$GOPATH/src/github.com/cosmos/cosmos-sdk"
SDK_BRANCH=v0.40.0-rc0
SIMD_DATA="$(pwd)/data"

# ARGS: 
# $1 -> local || remote, defaults to remote

# Ensure user understands what will be deleted
if [[ -d $SIMD_DATA ]] && [[ ! "$2" == "skip" ]]; then
  read -p "$0 will delete \$(pwd)/data folder. Do you wish to continue? (y/n): " -n 1 -r
  echo 
  if [[ ! $REPLY =~ ^[Yy]$ ]]; then
      exit 1
  fi
fi

rm -rf $SIMD_DATA #&> /dev/null
killall simd #&> /dev/null

set -e


if [[ -d $SDK_REPO ]]; then
  cd $SDK_REPO

  # remote build syncs with remote then builds
  if [[ "$1" == "local" ]]; then
    echo "Using local version of github.com/cosmos/cosmos-sdk"
    make simd &> /dev/null
    cp ./build/simd $GOBIN
  else
    echo "Building github.com/cosmos/cosmos-sdk@$SDK_BRANCH..."
    if [[ ! -n $(git status -s) ]]; then
      # sync with remote $SDK_BRANCH
      git fetch --all #&> /dev/null

      # ensure the gaia repository successfully pulls the latest $SDK_BRANCH
      if [[ -n $(git checkout $SDK_BRANCH -q) ]] || [[ -n $(git pull origin $SDK_BRANCH -q) ]]; then
        echo "failed to sync remote branch $SDK_BRANCH"
        echo "in $SDK_REPO, please rename the remote repository github.com/cosmos/cosmos-sdk to 'origin'"
        exit 1
      fi

      # install
      make simd &> /dev/null
      cp ./build/simd $GOBIN

      # ensure that built binary has the same version as the repo
      if [[ ! "$(simd version --long 2>&1 | grep "commit:" | sed 's/commit: //g')" == "$(git rev-parse HEAD)" ]]; then
        echo "built version of simd commit doesn't match"
        exit 1
      fi 
    else
      echo "uncommited changes in $SDK_REPO, please commit or stash before building"
      exit 1
    fi
    
  fi 
else 
  echo "$SDK_REPO doesn't exist, and you may not have have the gaia repo locally,"
  echo "if you want to download gaia to your \$GOPATH try running the following command:"
  echo "mkdir -p $(dirname $SDK_REPO) && git clone git@github.com:cosmos/gaia $SDK_REPO"
fi