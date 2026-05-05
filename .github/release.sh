#!/bin/bash

# Run before running this script: `git flow start release/<version>`

type=$1    # release o hotfix
version=$2 # es. 1.2.0 o 1.2.1-hotfix

if [ -z "$type" ] || [ -z "$version" ]; then
    echo "Usage: $0 <release|hotfix> <version> (e.g., release 1.2.0)"
    exit 1
fi

SOURCE_BRANCH="release/$version"
git switch $SOURCE_BRANCH
git flow release publish "$version"
git push
git flow release finish "$version" --nodevelopmerge -Fp
