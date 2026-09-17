#!/bin/sh
# ==============================================================================
# TenPoint — entrypoint cho container `web` (Next.js)
#
# Vì sao cần script này:
#   frontend/next.config.ts hiện CHƯA bật `output: "standalone"`.
#   - Nếu ĐÃ bật  → build sinh ra .next/standalone/server.js → chạy `node server.js`
#                   (image nhỏ, khởi động nhanh, không cần node_modules đầy đủ).
#   - Nếu CHƯA bật → fallback sang `npm start` (next start) với node_modules đầy đủ.
#
#   Fallback chỉ để container vẫn chạy được trong lúc FE chưa kịp bật standalone.
#   Xem mục "Việc cần FE làm" trong README.md ở thư mục gốc.
# ==============================================================================
set -eu

export HOSTNAME="${HOSTNAME:-0.0.0.0}"
export PORT="${PORT:-3000}"

if [ -f /app/server.js ]; then
    echo "[web] chế độ standalone — node server.js (PORT=${PORT} HOSTNAME=${HOSTNAME})"
    exec node /app/server.js
fi

echo "[web] ⚠️  KHÔNG tìm thấy /app/server.js — next.config.ts chưa bật output: 'standalone'."
echo "[web] ⚠️  Đang fallback sang 'npm start'. Image nặng hơn và khởi động chậm hơn."
echo "[web] ⚠️  Đề nghị FE thêm  output: 'standalone'  vào frontend/next.config.ts."
exec npm run start -- --port "${PORT}" --hostname "${HOSTNAME}"
