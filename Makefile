# D&D Simulator Makefile

.PHONY: help build run dev-up dev-down prod-up prod-down clean test

help: ## Show this help message
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

build: ## Build the Go application
	go build -o dnd-simulator .

run: ## Run the application locally (requires local MongoDB)
	go run main.go

dev-up: ## Start development environment (MongoDB + Mongo Express)
	docker-compose -f docker-compose.dev.yml up -d
	@echo "MongoDB running on localhost:27017"
	@echo "Mongo Express UI running on localhost:8081"

dev-down: ## Stop development environment
	docker-compose -f docker-compose.dev.yml down

prod-up: ## Start production environment (all services)
	docker-compose up -d --build
	@echo "Application running on localhost:8080"
	@echo "MongoDB running on localhost:27017"

prod-down: ## Stop production environment
	docker-compose down

prod-logs: ## View production logs
	docker-compose logs -f

clean: ## Clean up Docker containers and volumes
	docker-compose down -v
	docker-compose -f docker-compose.dev.yml down -v
	docker system prune -f

test: ## Run tests
	go test ./...

deps: ## Download Go dependencies
	go mod download
	go mod tidy

# Development workflow commands
dev-run: dev-up ## Start dev environment and run app locally
	@echo "Waiting for MongoDB to be ready..."
	@sleep 5
	MONGO_URI=mongodb://admin:password123@localhost:27017/dnd_simulator_dev?authSource=admin go run main.go

docker-build: ## Build Docker image
	docker build -t dnd-simulator .

# Database commands
db-shell: ## Connect to MongoDB shell
	docker exec -it dnd-mongodb-dev mongosh mongodb://admin:password123@localhost:27017/dnd_simulator_dev?authSource=admin