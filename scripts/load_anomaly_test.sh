#!/bin/bash
set -e

HOST=${1:-http://localhost:8081}
NORMAL_POINTS=${2:-50}
NORMAL_RPS=${3:-100}
ANOMALY_RPS=${4:-2000}

echo "== Load test started =="
echo "Target: $HOST"
echo "Normal points: $NORMAL_POINTS"
echo "Normal RPS: $NORMAL_RPS"
echo "Anomaly RPS: $ANOMALY_RPS"
echo

# WARMUP 
echo "[1/3] Warming up rolling window..."
for ((i=1; i<=NORMAL_POINTS; i++)); do
  curl -s -X POST "$HOST/metrics" \
    -H "Content-Type: application/json" \
    -d "{\"timestamp\":$i,\"cpu\":20,\"rps\":$NORMAL_RPS}" \
    >/dev/null || echo "WARN: request $i failed"
  sleep 0.02
done

# ANOMALY
echo
echo "[2/3] Sending anomaly point..."
curl -s -X POST "$HOST/metrics" \
  -H "Content-Type: application/json" \
  -d "{\"timestamp\":999999,\"cpu\":95,\"rps\":$ANOMALY_RPS}" \
  >/dev/null || echo "WARN: anomaly request failed"

# ANALYZE
echo
echo "[3/3] Fetching analysis..."
curl -s "$HOST/analyze"
echo
echo "== Load test finished =="
