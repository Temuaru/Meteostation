
##@ Общие команды:
.PHONY: help
help: ## Список всех доступных команд
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)


.PHONY: frontend-install 
frontend-install: ## Установка зависимостей node_modules
	@echo "--- Установка зависимостей node_modules ---"
	npm install

.PHONY: frontend-build
frontend-build: ## Сборка Webpack
	@echo "--- Сборка Webpack ---"
	npm run build

.PHONY: watch	
watch: ## Запуск режима наблюдения Webpack
	@echo "--- Запуск режима наблюдения Webpack ---"
	npm run watch

.PHONY: build	
build: ## Сборка контейнеров Docker
	@echo "--- Сборка контейнеров Docker ---"
	docker-compose stop
	docker-compose build --no-cache

.PHONY: down
down: ## Остановка и очистка контейнеров Docker
	@echo "--- Остановка и очистка контейнеров Docker ---"
	docker-compose stop
	docker-compose down --volumes

.PHONY: stop
stop: ## Остановка контейнеров Docker
	@echo "--- Остановка контейнеров Docker ---"
	docker-compose stop

.PHONY: run
run: ## Запуск контейнеров Docker
	@echo "--- Запуск контейнеров Docker ---"
	docker-compose up