#!/bin/bash

# Run before running this script: `git flow start release/<version>`

type=$1       # release o hotfix
version=$2    # es. 1.2.0 o 1.2.1-hotfix


if [ -z "$type" ] || [ -z "$version" ]; then
  echo "Usage: $0 <release|hotfix> <version> (e.g., release 1.2.0)"
  exit 1
fi

# Ensure local tag list is up-to-date with remote
git fetch --tags >/dev/null 2>&1 || true

# Check if the tag already exists locally or (after fetch) remotely.
# Try both the exact version and a common "v"-prefixed form.
if git rev-parse -q --verify "refs/tags/$version" >/dev/null 2>&1 || \
   git rev-parse -q --verify "refs/tags/v$version" >/dev/null 2>&1; then
  echo "ERROR: Git tag '$version' (or 'v$version') already exists. Aborting."
  exit 1
fi


SOURCE_BRANCH="release/$version"
git switch $SOURCE_BRANCH
git flow release publish "$version"
git push --tags
git flow release finish "$version" --nodevelopmerge -Fp

git switch main
