#!/bin/sh
set -eu

cd /app

if [ ! -f go.mod ]; then
  echo "apps/backend/go.mod が見つかりません。"
  echo "先にホスト側で Go モジュールを初期化してください。"
  echo "例: cd apps/backend && go mod init github.com/rigarashi1024/sns_only_event_saimple/apps/backend"
  sleep infinity
fi

if [ ! -f cmd/api/main.go ]; then
  echo "apps/backend/cmd/api/main.go が見つかりません。"
  echo "先に最小の API エントリポイントを作成してください。"
  sleep infinity
fi

go mod tidy

exec air -c .air.toml
