# deploy.makefile

KUBECTL := kubectl
NAMESPACE ?= gobackend
CONTEXT ?= sre.im

DEPLOYS = gobackend-apiserver

.PHONY: deploy.run.all
deploy.run.all:
	@echo "==========> Deploying all"
	@$(MAKE) deploy.run

.PHONY: deploy.run
deploy.run: $(addprefix deploy.run., ${DEPLOYS})

.PHONY: deploy.run.%
deploy.run.%:
	$(eval ARCH := $(word 2,$(subst _, ,${PLATFORM})))
	@echo "==========> Deploying $* ${VERSION} ${ARCH}"
	echo @${KUBECTL} -n ${NAMESPACE} --context=${CONTEXT} set image deployments/$* $*=${REGISTRY_PREFIX}/$*:${VERSION}

