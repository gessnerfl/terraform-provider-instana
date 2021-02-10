#!/bin/sh
github_changelog_generator --token $2 -u $1 -p terraform-provider-instana --enhancement-labels "improvement,feature,enhancement"