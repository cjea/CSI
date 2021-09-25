.PHONY: phony
run: phony ## Run the app.
	@ go run main.go https://www.ebay.com

test: phony ## Execute unit tests.
	go run github.com/onsi/ginkgo/ginkgo -r

BLUE := $(shell tput setaf 4)
RESET := $(shell tput sgr0)

.PHONY: help
help: ## List all targets and short descriptions of each.
	@grep -E '^[^ .]+: .*?## .*$$' $(MAKEFILE_LIST) \
		| sort \
		| awk '\
			BEGIN { FS = ": .*##" };\
			{ printf "$(BLUE)%-29s$(RESET) %s\n", $$1, $$2  }'
