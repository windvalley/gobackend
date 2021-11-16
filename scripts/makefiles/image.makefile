# image.makefile

DOCKER := docker
DOCKER_SUPPORTED_API_VERSION ?= 1.32

BASE_IMAGE = alpine:3
REGISTRY_PREFIX ?= private-registry.xxx.com/windvalley

EXTRA_ARGS ?= --no-cache
_DOCKER_BUILD_EXTRA_ARGS :=

ifdef HTTP_PROXY
	_DOCKER_BUILD_EXTRA_ARGS += --build-arg HTTP_PROXY=${HTTP_PROXY}
endif

ifneq ($(EXTRA_ARGS), )
	_DOCKER_BUILD_EXTRA_ARGS += $(EXTRA_ARGS)
endif

IMAGES_DIR ?= $(wildcard ${ROOT_DIR}/build/docker/*)
IMAGES ?= $(filter-out tools,$(foreach image,${IMAGES_DIR},$(notdir ${image})))

ifeq (${IMAGES},)
	$(error Could not determine IMAGES, set ROOT_DIR or run in source dir)
endif

.PHONY: image.verify
image.verify:
	$(eval API_VERSION := $(shell $(DOCKER) version | grep -E 'API version: {1,6}[0-9]' | head -n1 | awk '{print $$3} END { if (NR==0) print 0}' ))
	$(eval PASS := $(shell echo "$(API_VERSION) > $(DOCKER_SUPPORTED_API_VERSION)" | bc))
	@if [ $(PASS) -ne 1 ]; then \
		$(DOCKER) -v ;\
		echo "Unsupported docker version. Docker API version should be greater than $(DOCKER_SUPPORTED_API_VERSION)"; \
		exit 1; \
	fi

.PHONY: image.build
image.build: image.verify go.build.verify $(addprefix image.build., $(addprefix ${IMAGE_PLAT}., ${IMAGES}))

.PHONY: image.build.%
image.build.%: go.build.%
	$(eval IMAGE := ${COMMAND})
	$(eval IMAGE_PLAT := $(subst _,/,${PLATFORM}))
	@echo "==========> Building docker image ${IMAGE} ${VERSION} for ${IMAGE_PLAT}"
	@mkdir -p ${TMP_DIR}/${IMAGE}
	@cat ${ROOT_DIR}/build/docker/${IMAGE}/Dockerfile\
		| sed "s#BASE_IMAGE#${BASE_IMAGE}#" >${TMP_DIR}/${IMAGE}/Dockerfile
	@cp ${OUTPUT_DIR}/platforms/${IMAGE_PLAT}/${IMAGE} ${TMP_DIR}/${IMAGE}/
	@cp configs/*${COMMAND}.yaml ${TMP_DIR}/${IMAGE}/
	$(eval BUILD_SUFFIX := ${_DOCKER_BUILD_EXTRA_ARGS} --pull -t ${REGISTRY_PREFIX}/${IMAGE}:${VERSION} ${TMP_DIR}/${IMAGE})
	@${DOCKER} build ${BUILD_SUFFIX}


