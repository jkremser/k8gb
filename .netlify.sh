#!/bin/bash

if [[ $CI == "true" ]]; then
    echo "Preparing the preview of website for deployment"
    if ! git config remote.origin.url &> /dev/null; then
        git remote add -f -t gh-pages origin https://github.com/k8gb-io/k8gb
    fi
    git fetch origin gh-pages:gh-pages
    git checkout gh-pages
    git checkout - {README,CONTRIBUTING,CHANGELOG}.md docs/
    mv CNAME EMANC
    bundle install
    bundle exec jekyll build
    mv EMANC CNAME
    git checkout -
else
    echo "This script is no-op when running outside CI"
fi
