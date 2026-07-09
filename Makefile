APP_NAME=new-rag-go-app-service
APP_CMD=./cmd/web

.PHONY: up up_build down build logs clean restart shell

up:
	@echo "Starting Docker services..."
	APP_CMD=${APP_CMD} docker compose up
	@echo "Docker services started!"

up_build:
	@echo "Stopping Docker services if running..."
	docker compose down
	@echo "Building and starting Docker services..."
	APP_CMD=${APP_CMD} docker compose up --build
	@echo "Docker services started!"

build:
	@echo "Building Docker images..."
	APP_CMD=${APP_CMD} docker compose build
	@echo "Done!"

down:
	@echo "Stopping Docker services..."
	docker compose down

restart: down up_build

logs:
	docker compose logs -f

shell:
	docker exec -it ${APP_NAME} sh

clean:
	@echo "Removing containers and volumes..."
	docker compose down -v
	@echo "Done!"