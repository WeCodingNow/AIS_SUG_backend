BACKEND_PATH=.
SERVER_PATH=$(BACKEND_PATH)/cmd/server/main.go
BACKEND_PORT=8080

BACKEND_IMAGE=ais-sug-apiserver
TAG=ver1

.PHONY: build run

image:
	docker build -t $(BACKEND_IMAGE):$(TAG) $(BACKEND_PATH)

run:
	go run $(SERVER_PATH)

docker: image
	docker run \
		-p $(BACKEND_PORT):8080 \
		$(BACKEND_IMAGE):$(TAG) 
