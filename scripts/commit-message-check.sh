#!/bin/bash

commit_regex='^(feat|fix|docs|chore): .+'

if ! grep -qE "$commit_regex" "$1"; then
    echo "❌ Invalid commit message format!"
    echo ""
    echo "Format: 'type: description'"
    echo "Types: feat, fix, docs, chore"
    echo "Example: 'feat: add user authentication'"
    exit 1
fi