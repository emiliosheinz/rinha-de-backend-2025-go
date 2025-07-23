#!/bin/bash

docker compose down --volumes --remove-orphans
docker compose up -d --build --force-recreate

sleep 3

docker run --rm -i \
  --network host \
  -v $(pwd)/rinha-test:/scripts \
  -w /scripts \
  grafana/k6 run rinha.js
