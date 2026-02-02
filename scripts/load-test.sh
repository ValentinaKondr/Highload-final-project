#!/bin/bash

set -e

SERVICE_URL=${1:-"http://localhost:8080"}
RPS=${2:-1000}
DURATION=${3:-300}

echo "Starting load test..."
echo "Service URL: $SERVICE_URL"
echo "Target RPS: $RPS"
echo "Duration: ${DURATION}s"

# Check if locust is installed
if ! command -v locust &> /dev/null; then
    echo "Locust not found. Installing..."
    pip3 install locust
fi

# Create locustfile if it doesn't exist
cat > locustfile.py << 'EOF'
from locust import HttpUser, task, between
import random
import time

class MetricsUser(HttpUser):
    wait_time = between(0.1, 0.5)
    
    @task(3)
    def send_metric(self):
        timestamp = int(time.time())
        cpu = random.uniform(10, 90)
        rps = random.uniform(100, 2000)
        
        payload = {
            "timestamp": timestamp,
            "cpu": cpu,
            "rps": rps
        }
        self.client.post("/metrics", json=payload, name="Send Metric")
    
    @task(1)
    def get_analyze(self):
        self.client.get("/analyze", name="Get Analytics")
    
    @task(1)
    def get_health(self):
        self.client.get("/health", name="Health Check")
EOF

echo "Running Locust..."
locust -f locustfile.py --headless -u $RPS -r 100 --run-time ${DURATION}s --host=$SERVICE_URL --html=xtemp/load-test-report.html

echo "Load test completed! Report saved to xtemp/load-test-report.html"

