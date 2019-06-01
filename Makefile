
.PHONY: build-doke
build-doke: ## build doke
	go build .

.PHONY: print-an-env-var
print-an-env-var:
	@echo ${MY_ENV_VAR}


.PHONY: help
help:
	@grep -E -h '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help
