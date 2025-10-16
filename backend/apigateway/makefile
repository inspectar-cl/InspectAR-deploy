APP_NAME=apigateway
CMD_DIR=./cmd

.PHONY: tidy run build clean

tidy:
	@echo ">>> Verificando dependencias..."
	@go mod tidy

run: tidy
	@echo ">>> Ejecutando $(APP_NAME) ..."
	@go run $(CMD_DIR)/main.go

build: tidy
	@echo ">>> Compilando binario $(APP_NAME) ..."
	@go build -o bin/$(APP_NAME) $(CMD_DIR)/main.go

clean:
	@echo ">>> Limpiando binarios ..."
	@rm -rf bin

