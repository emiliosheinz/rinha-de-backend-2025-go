#!/bin/bash

docker compose down --volumes --remove-orphans
docker compose up -d --build --force-recreate

sleep 10

K6_WEB_DASHBOARD=true K6_WEB_DASHBOARD_OPEN=true k6 run ./rinha-test/rinha.js
