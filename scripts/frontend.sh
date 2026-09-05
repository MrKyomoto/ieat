#!/usr/bin/env bash
set -e

project_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$project_root/web"

exec npm run dev
