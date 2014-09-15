BUILD_DIR := image
GIT_BRANCH := $(shell git rev-parse --abbrev-ref HEAD 2>/dev/null)
GITCOMMIT := $(shell git rev-parse --short HEAD 2>/dev/null)
PROJECT_NAME := pty
RELEASE_IMAGE := deis/pty
BUILD_IMAGE := $(PROJECT_NAME)-build$(if $(GIT_BRANCH),:$(GIT_BRANCH))
DOCKER_IMAGE := $(PROJECT_NAME)$(if $(GIT_BRANCH),:$(GIT_BRANCH))

build:
	docker build -t $(BUILD_IMAGE) .

image: build
	docker run -v ${PWD}/$(BUILD_DIR):/$(BUILD_DIR):rw $(BUILD_IMAGE) cp /go/bin/pty /$(BUILD_DIR)/
	docker build -t $(DOCKER_IMAGE) $(BUILD_DIR)
	rm -f $(BUILD_DIR)/pty

release: image
	docker tag $(DOCKER_IMAGE) $(RELEASE_IMAGE)
	docker push $(RELEASE_IMAGE)

clean:
	docker rmi $(BUILD_IMAGE) $(DOCKER_IMAGE)
