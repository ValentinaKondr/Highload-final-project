from locust import HttpUser, task, between
import random
import time

class MetricsUser(HttpUser):
    """
    Locust user class for load testing the metrics service.
    Simulates IoT devices sending metrics.
    """
    wait_time = between(0.1, 0.5)
    
    @task(3)
    def send_metric(self):
        """Send a metric (most common task)"""
        timestamp = int(time.time())
        cpu = random.uniform(10, 90)
        rps = random.uniform(100, 2000)
        
        payload = {
            "timestamp": timestamp,
            "cpu": cpu,
            "rps": rps
        }
        with self.client.post("/metrics", json=payload, name="Send Metric", catch_response=True) as response:
            if response.status_code == 200:
                response.success()
            else:
                response.failure(f"Got status code {response.status_code}")
    
    @task(1)
    def get_analyze(self):
        """Get analytics data"""
        self.client.get("/analyze", name="Get Analytics")
    
    @task(1)
    def get_health(self):
        """Health check"""
        self.client.get("/health", name="Health Check")
    
    @task(1)
    def send_anomaly(self):
        """Occasionally send anomalous values"""
        timestamp = int(time.time())
        cpu = random.uniform(10, 90)
        # Anomalous RPS value (very high)
        rps = random.uniform(5000, 10000)
        
        payload = {
            "timestamp": timestamp,
            "cpu": cpu,
            "rps": rps
        }
        self.client.post("/metrics", json=payload, name="Send Anomaly")

