URL="http://localhost:9999/payments"

for i in {1..3}; do
  amount=$((RANDOM % 1000 + 1))
  user_id=$((RANDOM % 5 + 1))
  correlation_id=$(uuidgen)

  curl -s -X POST "$URL" \
    -H "Content-Type: application/json" \
    -d "{\"amount\": $amount, \"correlationId\": \"$correlation_id\"}" &
done

wait
echo "✅ Requisições enviadas para $URL com valores e correlationId aleatórios"
