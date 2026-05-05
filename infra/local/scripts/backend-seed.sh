#!/bin/sh
set -eu

cd /app

if [ ! -f go.mod ]; then
  echo "apps/backend/go.mod が見つかりません。"
  exit 1
fi

if [ ! -f cmd/seed/main.go ]; then
  echo "apps/backend/cmd/seed/main.go が見つかりません。"
  exit 1
fi

max_attempts="${SEED_MAX_ATTEMPTS:-30}"
sleep_seconds="${SEED_RETRY_SECONDS:-2}"
attempt=1

while [ "$attempt" -le "$max_attempts" ]; do
  echo "Firestore seed を実行します... (${attempt}/${max_attempts})"
  if go run ./cmd/seed; then
    echo "Firestore seed が完了しました。"
    exit 0
  fi

  echo "Firestore seed に失敗しました。${sleep_seconds} 秒後に再試行します。"
  attempt=$((attempt + 1))
  sleep "$sleep_seconds"
done

echo "Firestore seed が最大試行回数を超えて失敗しました。"
exit 1
