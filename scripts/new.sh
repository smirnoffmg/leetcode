#!/usr/bin/env bash
set -euo pipefail

if [ $# -lt 1 ]; then
	echo "usage: $0 <slug> [problem-number] [title] [difficulty]" >&2
	echo "example: $0 two-sum" >&2
	echo "example: $0 two-sum 1 \"Two Sum\" easy   # фолбэк, если leetcode недоступен" >&2
	exit 1
fi

slug=$1
num=${2:-0}
title=${3:-}
difficulty=${4:-}

root=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)

# Номер, название и сложность нужны только как фолбэк: если leetcode ответил,
# всё берётся из его ответа.
exec go run "$root/cmd/lcgen" \
	-slug "$slug" \
	-num "$num" \
	-title "$title" \
	-difficulty "$difficulty" \
	-root "$root"
