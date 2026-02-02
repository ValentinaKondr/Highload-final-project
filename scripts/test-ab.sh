#!/bin/bash

set -e

SERVICE_URL=${1:-"http://localhost:8080"}
TOTAL_REQUESTS=${2:-50000}
CONCURRENCY=${3:-100}

echo "Starting Apache Bench load test..."
echo "Service URL: $SERVICE_URL"
echo "Total requests: $TOTAL_REQUESTS"
echo "Concurrency: $CONCURRENCY"

# Check if ab is installed
if ! command -v ab &> /dev/null; then
    echo "Apache Bench not found. Installing..."
    if [[ "$OSTYPE" == "darwin"* ]]; then
        brew install httpd
    elif [[ "$OSTYPE" == "linux-gnu"* ]]; then
        sudo apt-get update && sudo apt-get install -y apache2-utils
    fi
fi

# Create test data file
cat > /tmp/metrics.json << 'EOF'
{"timestamp": 1234567890, "cpu": 45.5, "rps": 1500}
EOF

# Run Apache Bench
echo "Running Apache Bench..."
ab -n $TOTAL_REQUESTS -c $CONCURRENCY -p /tmp/metrics.json -T application/json \
  -g xtemp/ab-results.tsv \
  $SERVICE_URL/metrics > xtemp/ab-output.txt

echo "Test completed!"
echo "Results saved to:"
echo "  - xtemp/ab-output.txt (summary)"
echo "  - xtemp/ab-results.tsv (detailed data)"

# Display summary
echo ""
echo "=== Summary ==="
grep -E "(Requests per second|Time per request|Failed requests)" xtemp/ab-output.txt

