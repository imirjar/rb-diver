include .env

.PHONY: run build test prod docker-build docker-run docker-clean



#docs:
#	swag init -g internal/gateway/http/http.go 

# Проходим все тесты
test: 
	go test ./...

# Запускаем в консоли
run:
	go run $(APP_FILE)

# Собираем бинарник
build:
	CGO_ENABLED=0 GOOS=linux go build -o bin/$(BIN_NAME) $(APP_FILE)

# Собираем бинарник запуск в консоли
start: build 
	./bin/$(BIN_NAME) &

# Docker сборка
docker-build: 
	docker build -t $(DOCKER_IMAGE_NAME) .

# Запуск проекта в Docker контейнере с использованием .env файла
docker-run: docker-build
	docker run -d --env-file .env --name $(DOCKER_CONTAINER_NAME) -p $(PORT):$(PORT) $(DOCKER_IMAGE_NAME)

# Очищение контейнера после остановки
docker-clean:
	@if [ "$(shell docker ps -a -q --filter "name=$(DOCKER_CONTAINER_NAME)")" ]; then \
		docker rm $(DOCKER_CONTAINER_NAME); \
	fi

# Полный цикл сборки и запуска в Docker
prod: docker-clean docker-run
# prod: docker-run
