A Go-based HTTP service that fetches data from multiple URLs in parallel and returns their HTTP status codes. This project is containerized with Docker, deployed to Kubernetes using Helm, and includes a GitHub Actions CI/CD pipeline for linting, testing, and building.

## Project Overview
The service exposes two endpoints:
- **POST /fetch**: Accepts a JSON payload with a list of URLs and returns their HTTP status codes or errors.
  ```json
  {"urls": ["https://httpbin.org/get"]}
  ```
  Example response:
  ```json
  {"results":{"https://httpbin.org/get":200},"errors":{}}
  ```
- **GET /health**: Returns "OK" for liveness/readiness probes.

## Prerequisites
- Go 1.24.3
- Docker
- Minikube (for local Kubernetes)
- Helm
- Git
- curl or Postman (for testing)

## Setup Instructions
1. **Clone the Repository**:
   ```bash
   git clone https://github.com/your-username/parallel-data-fetcher.git
   cd parallel-data-fetcher
   ```

2. **Build the Docker Image**:
   ```bash
   docker build -t parallel-fetcher:latest .
   ```

3. **Run Locally**:
   ```bash
   docker run -p 8080:8080 parallel-fetcher:latest
   curl -X POST "http://localhost:8080/fetch" -H "Content-Type: application/json" -d '{"urls": ["https://httpbin.org/get"]}'
   ```

4. **Deploy to Minikube**:
   ```bash
   minikube start --driver=docker
   minikube image load parallel-fetcher:latest
   helm install fetcher ./parallel-fetcher
   kubectl port-forward svc/fetcher-parallel-fetcher 8083:80
   curl -X POST "http://localhost:8083/fetch" -H "Content-Type: application/json" -d '{"urls": ["https://httpbin.org/get"]}'
   ```

## CI/CD Pipeline
A GitHub Actions workflow (`.github/workflows/ci.yml`) automates the following on push or pull requests to the `master` branch:
- **Linting**: Runs `go vet` to ensure code quality.
- **Testing**: Executes unit tests with `go test ./...` (see `main_test.go`).
- **Building**: Builds the Docker image with `docker build -t parallel-fetcher:latest .`.
- **Optional Deployment**: Can deploy to a test Kubernetes namespace (`test-fetcher`) using Helm (commented out in `ci.yml`).

View CI runs in the GitHub repository's Actions tab.

## Project Structure
```
├── Dockerfile
├── go.mod
├── main.go
├── main_test.go
├── parallel-fetcher/
│   ├── Chart.yaml
│   ├── values.yaml
│   └── templates/
│       ├── deployment.yaml
│       ├── service.yaml
│       └── NOTES.txt
├── .github/
│   └── workflows/
│       └── ci.yml
└── README.md
```

## Testing
Run unit tests locally:
```bash
go test ./... -v
```

Test the service:
```bash
curl -X POST "http://localhost:8080/fetch" -H "Content-Type: application/json" -d '{"urls": ["https://httpbin.org/get"]}'
```

## Contributing
Feel free to open issues or submit pull requests for improvements.

## License
MIT License