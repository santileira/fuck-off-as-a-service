SHELL=/bin/sh

.PHONY: fmt
fmt:
	@echo "----------------------------------------------------------------"
	@echo " ⚙️  Formatting code..."
	@echo "----------------------------------------------------------------"
	go fmt ./...
	go mod tidy

.PHONY: test
test:
	@echo "----------------------------------------------------------------"
	@echo " ✅  Testing code..."
	@echo "----------------------------------------------------------------"
	go test ./... -coverprofile=coverage.out

.PHONY: coverage
coverage:
	@echo "----------------------------------------------------------------"
	@echo " 📊  Checking coverage..."
	@echo "----------------------------------------------------------------"
	go tool cover -html=coverage.out -o coverage.html
	go tool cover -func=coverage.out

.PHONY: deps
deps:
	@echo "----------------------------------------------------------------"
	@echo " ⬇️  Downloading dependencies..."
	@echo "----------------------------------------------------------------"
	go get ./...

.PHONY: build
build: deps fmt
	@echo "----------------------------------------------------------------"
	@echo " 📦 Building binary..."
	@echo "----------------------------------------------------------------"
	go build -o fuck-off-as-a-service main.go

.PHONY: run
run:
	@echo "----------------------------------------------------------------"
	@echo " ️🏃 Running..."
	@echo "----------------------------------------------------------------"
	./fuck-off-as-a-service serve

.PHONY: all
all: test build

.PHONY: docker-build
docker-build:
	@echo "----------------------------------------------------------------"
	@echo " 📦 ️Building in docker..."
	@echo "----------------------------------------------------------------"
	docker build -t fuck-off-as-a-service .

.PHONY: docker-run
docker-run:
	@echo "----------------------------------------------------------------"
	@echo " ️🏃 Running in docker..."
	@echo "----------------------------------------------------------------"
	docker run -p 8080:8080 -it --rm fuck-off-as-a-service