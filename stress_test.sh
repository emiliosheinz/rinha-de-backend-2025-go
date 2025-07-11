#!/bin/bash

URL="http://localhost:9999/payments"

for i in {1..1000}; do
  curl -s -X POST "$URL" -H "Content-Type: application/json" -d '{"amount": 100, "user_id": 1}' &
done

wait
echo "✅ 1000 requisições enviadas para $URL"
