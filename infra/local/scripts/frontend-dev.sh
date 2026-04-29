#!/bin/sh
set -eu

cd /app

if [ ! -f package.json ]; then
  echo "apps/frontend/package.json が見つかりません。"
  echo "先にホスト側で Nuxt を初期化してください。"
  echo "例: cd apps/frontend && npm create nuxt@latest ."
  sleep infinity
fi

if [ ! -d node_modules ] || [ -z "$(ls -A node_modules 2>/dev/null)" ]; then
  echo "frontend dependencies をインストールします..."
  npm install
fi

exec npm run dev -- --host 0.0.0.0 --port 3001
