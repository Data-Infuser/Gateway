export APP=gateway
export CONTAINER=infuser-gateway
export CONTAINER_VERSION=0.1
export CONTAINER_PORT=9090
export EXPOSE_PORT=9090

build:
	go build ./main.go

container:
	docker build --tag $(CONTAINER):$(CONTAINER_VERSION) .

run-container:
	docker run --rm --detach --publish $(EXPOSE_PORT):$(CONTAINER_PORT) --name $(APP) $(CONTAINER):$(CONTAINER_VERSION)

container-log:
	docker logs --follow $(APP)