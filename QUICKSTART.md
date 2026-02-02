# Быстрый старт

## Минимальные требования

- Go 1.22+
- Docker
- Kubernetes (Minikube или Kind)
- kubectl

## Развертывание

### 1. Запуск Kubernetes кластера

**Minikube:**
```bash
minikube start --cpus=2 --memory=4g
minikube addons enable metrics-server
minikube addons enable ingress
```

**Kind (опционально):**
```bash
kind create cluster --config kind-config.yaml
```

### 2. Сборка и загрузка образа

```bash
# Сборка образа
./scripts/build.sh

# Для Minikube загрузите образ в локальный реестр:
eval $(minikube docker-env)
docker build -t go-service:latest .
```

### 3. Развертывание

```bash
# Развертывание всех компонентов
./scripts/deploy.sh

# Проверка статуса
kubectl get pods,svc,hpa
```
#### Будут развернуты:

- Go-сервис
- Redis
- ConfigMap и Secret
- Horizontal Pod Autoscaler (HPA)

### 4. Доступ к сервису

```bash
# Port-forward
kubectl port-forward service/go-service 8080:80

# Тестирование
curl http://localhost:8080/health
```

### 5. Мониторинг (опционально)

```bash
# Установка Prometheus и Grafana
./scripts/setup-monitoring.sh

# Доступ к интерфейсам
kubectl port-forward -n monitoring svc/prometheus-kube-prometheus-prometheus 9090:9090
kubectl port-forward -n monitoring svc/prometheus-grafana 3000:80
# Логин: admin / admin123
```

### 6. Нагрузочное тестирование

```bash
# Установка Locust
pip3 install locust

# Запуск теста
kubectl port-forward service/go-service 8080:80 &
./scripts/load-test.sh http://localhost:8080 1000 300
```

Еще один скрипт нагрузки
```bash
./scripts/load_anomaly_test.sh http://localhost:8080
```
Apache Bench (для тестирования HPA)
```bash
ab -n 100000 -c 200 \
  -p payload.json \
  -T application/json \
  http://localhost:8080/metrics
````
мониторинг масштабирования
```bash
watch kubectl get hpa,pods
```
## Проверка работы

### Health Check
```bash
curl http://localhost:8080/health
```

### Отправка метрики
```bash
curl -X POST http://localhost:8080/metrics \
  -H "Content-Type: application/json" \
  -d '{"timestamp": '$(date +%s)', "cpu": 45.5, "rps": 1500}'
```

### Получение аналитики
```bash
curl http://localhost:8080/analyze
```

### Prometheus метрики
```bash
curl http://localhost:8080/metrics
```

## Мониторинг масштабирования

```bash
# Наблюдение за HPA и подами
watch kubectl get hpa,pods

# Генерация нагрузки для тестирования масштабирования
for i in {1..5}; do
  kubectl run load-gen-$i --image=busybox --rm -it --restart=Never -- \
    sh -c "while true; do wget -q -O- http://go-service/metrics; done" &
done
```

## Очистка

```bash
kubectl delete -f k8s/
minikube stop
```

## Дополнительная документация

- Полное руководство: `README.md`
- Развертывание: `xtemp/DEPLOYMENT_GUIDE.md`
- Тестирование: `xtemp/TESTING_GUIDE.md`
- Отчет: `xtemp/PROJECT_REPORT_TEMPLATE.md`

