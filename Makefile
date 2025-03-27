UNAME_S := $(shell uname -s)

# Detect the operating system and set the OS variable
ifeq ($(UNAME_S),Linux)
    OS := Linux
else ifeq ($(UNAME_S),Darwin)
    OS := MacOS
else
    OS := Unknown
endif

.PHONY: install_swag
install_swag:
	go install github.com/swaggo/swag/cmd/swag@latest
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
	@echo "Detected Operating System: $(OS)"
	@if [ "$(OS)" = "Linux" ]; then \
		npm config set prefix "~/.npm"; \
		npm install @openapitools/openapi-generator-cli -g; \
		ln -s ~/.npm/bin/openapi-generator-cli ~/.local/bin/openapi-generator; \
		openapi-generator-cli version-manager set 7.6.0; \
	elif [ "$(OS)" = "MacOS" ]; then \
		brew install openapi-generator; \
	else \
		echo "Unsupported Operating System: $(OS)"; \
		exit 1; \
	fi

.PHONY: install_mockery
install_mockery:
	go install github.com/vektra/mockery/v2@v2.42.0

.PHONY: generate_swagger
generate_swagger:
	swag init -g service/http/service.go -o cmd/service/http/docs -d ./cmd,./pkg/models --instanceName TraderBackend

.PHONY: generate_mocks
generate_mocks:
	mockery --name ExchangeClient --dir cmd/interfaces --output cmd/interfaces/mocks --filename exchange_client.go
	mockery --name Trader --dir cmd/interfaces --output cmd/interfaces/mocks --filename trader.go

.PHONY: docker
docker:
	docker build -t bybit-trader -f ./Dockerfile  .

.PHONY: test
test:
	go test ./...