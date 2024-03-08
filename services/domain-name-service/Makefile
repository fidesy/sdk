# Constants

PROJECT_NAME=domain-name-service
USER=fidesy-pay


PHONY: generate
generate:
	mkdir -p pkg/${PROJECT_NAME}
	protoc --go_out=pkg/${PROJECT_NAME} --go_opt=paths=import \
			--go-grpc_out=pkg/${PROJECT_NAME} --go-grpc_opt=paths=import \
			--grpc-gateway_out=pkg/${PROJECT_NAME} \
			--grpc-gateway_opt grpc_api_configuration=./api/${PROJECT_NAME}/${PROJECT_NAME}.yaml \
            --grpc-gateway_opt allow_delete_body=true \
			api/${PROJECT_NAME}/${PROJECT_NAME}.proto
	mv pkg/${PROJECT_NAME}/github.com/${USER}/${PROJECT_NAME}/* pkg/${PROJECT_NAME}
	rm -r pkg/${PROJECT_NAME}/github.com

PHONY: clean
clean:
	 if docker inspect ${PROJECT_NAME} > /dev/null 2>&1; then docker rm -f ${PROJECT_NAME} && docker rmi -f ${PROJECT_NAME}; else echo "Container not found."; fi

PHONY: go-build
go-build:
	GOOS=linux GOARCH=amd64 go build -o ./main ./cmd/${PROJECT_NAME}
	mkdir -p bin
	mv main bin

PHONY: build
build:
	make go-build
	docker build --tag ${PROJECT_NAME} .

PHONY: run
run:
	make clean

	make build

	docker run --name ${PROJECT_NAME} -d \
			--network zoo \
			-e APP_NAME=${PROJECT_NAME} \
			-e ENV=PRODUCTION \
			-e GRPC_PORT=10000 -p 10000:10000 \
			-e METRICS_PORT=10001 -p 10001:10001 \
			${PROJECT_NAME}