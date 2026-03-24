# 运维AI助手 - CI/CD配置

> 版本: v1.0
> 更新日期: 2026-03-24

---

## 1. GitHub Actions工作流

### 1.1 CI流程 (.github/workflows/ci.yml)

```yaml
name: CI

on:
  push:
    branches: [ develop, main ]
  pull_request:
    branches: [ develop, main ]

jobs:
  lint:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      - name: golangci-lint
        uses: golangci/golangci-lint-action@v3

  test:
    runs-on: ubuntu-latest
    services:
      postgres:
        image: postgres:14
        env:
          POSTGRES_PASSWORD: test
        options: >-
          --health-cmd pg_isready
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5
      redis:
        image: redis:7
        options: >-
          --health-cmd "redis-cli ping"
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      - name: Run tests
        run: make test
      - name: Upload coverage
        uses: codecov/codecov-action@v3

  build:
    runs-on: ubuntu-latest
    needs: [lint, test]
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      - name: Build
        run: make build
```

### 1.2 CD流程 (.github/workflows/cd.yml)

```yaml
name: CD

on:
  push:
    tags:
      - 'v*'

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3

      - name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v2

      - name: Login to Registry
        uses: docker/login-action@v2
        with:
          registry: ${{ secrets.REGISTRY_URL }}
          username: ${{ secrets.REGISTRY_USERNAME }}
          password: ${{ secrets.REGISTRY_PASSWORD }}

      - name: Build and push
        uses: docker/build-push-action@v4
        with:
          context: .
          file: ./deployments/docker/Dockerfile.api
          push: true
          tags: |
            ${{ secrets.REGISTRY_URL }}/ops-ai-api:${{ github.ref_name }}
            ${{ secrets.REGISTRY_URL }}/ops-ai-api:latest

      - name: Deploy to Kubernetes
        run: |
          kubectl set image deployment/ops-ai-api \
            ops-ai-api=${{ secrets.REGISTRY_URL }}/ops-ai-api:${{ github.ref_name }}
```

---

## 2. Dockerfile

### 2.1 API服务 (deployments/docker/Dockerfile.api)

```dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /api ./cmd/api

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /api .
EXPOSE 8080
CMD ["./api"]
```

### 2.2 Worker服务 (deployments/docker/Dockerfile.worker)

```dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /worker ./cmd/worker

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /worker .
CMD ["./worker"]
```

---

## 3. Kubernetes部署

### 3.1 API Deployment (deployments/kubernetes/api-deployment.yaml)

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: ops-ai-api
  namespace: ops-ai
spec:
  replicas: 3
  selector:
    matchLabels:
      app: ops-ai-api
  template:
    metadata:
      labels:
        app: ops-ai-api
    spec:
      containers:
      - name: ops-ai-api
        image: registry.example.com/ops-ai-api:latest
        ports:
        - containerPort: 8080
        env:
        - name: APP_ENV
          value: "production"
        envFrom:
        - configMapRef:
            name: ops-ai-config
        - secretRef:
            name: ops-ai-secrets
        resources:
          requests:
            cpu: 1000m
            memory: 2Gi
          limits:
            cpu: 2000m
            memory: 4Gi
        livenessProbe:
          httpGet:
            path: /healthz
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /readyz
            port: 8080
          initialDelaySeconds: 10
          periodSeconds: 5
```

### 3.2 Service (deployments/kubernetes/api-service.yaml)

```yaml
apiVersion: v1
kind: Service
metadata:
  name: ops-ai-api
  namespace: ops-ai
spec:
  selector:
    app: ops-ai-api
  ports:
  - protocol: TCP
    port: 80
    targetPort: 8080
  type: ClusterIP
```

---

## 4. 变更记录

| 版本 | 日期 | 变更 |
|------|------|------|
| v1.0 | 2026-03-24 | 初始版本 |
