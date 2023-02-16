# The version of the build harness container to use
BUILD_HARNESS_REPO := ghcr.io/defenseunicorns/not-a-build-harness/not-a-build-harness
BUILD_HARNESS_VERSION := 0.0.4

# The version of Zarf to use. To keep this repo as portable as possible the Zarf binary will be downloaded and added to
# the build folder.
ZARF_VERSION := v0.24.2

# Figure out which Zarf binary we should use based on the operating system we are on
ZARF_BIN := zarf
UNAME_S := $(shell uname -s)
UNAME_P := $(shell uname -p)
ifneq ($(UNAME_S),Linux)
	ifeq ($(UNAME_S),Darwin)
		ZARF_BIN := $(addsuffix -mac,$(ZARF_BIN))
	endif
	ifeq ($(UNAME_P),i386)
		ZARF_BIN := $(addsuffix -intel,$(ZARF_BIN))
	endif
	ifeq ($(UNAME_P),arm64)
		ZARF_BIN := $(addsuffix -apple,$(ZARF_BIN))
	endif
endif

.DEFAULT_GOAL := help

# Optionally add the "-it" flag for docker run commands if the env var "CI" is not set (meaning we are on a local machine and not in github actions)
TTY_ARG :=
ifndef CI
	TTY_ARG := -it
endif

# Silent mode by default. Run `make VERBOSE=1` to turn off silent mode.
ifndef VERBOSE
.SILENT:
endif

# Idiomatic way to force a target to always run, by having it depend on this dummy target
FORCE:

.PHONY: help
help: ## Show a list of all targets
	grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
	| sed -n 's/^\(.*\): \(.*\)##\(.*\)/\1:\3/p' \
	| column -t -s ":"

.PHONY: test
test: ## Run all automated tests. Requires access to an AWS account. Costs real money. Requires env vars REPO_URL, GIT_BRANCH, and AWS cred env vars
	mkdir -p .cache/go
	mkdir -p .cache/go-build
	echo "Running automated tests. This will take several minutes. At times it does not log anything to the console. If you interrupt the test run you will need to log into AWS console and manually delete any orphaned infrastructure."
	docker run $(TTY_ARG) --rm -v "${PWD}:/app" -v "${PWD}/.cache/go:/root/go" -v "${PWD}/.cache/go-build:/root/.cache/go-build" --workdir "/app/test/games_package" -e GOPATH=/root/go -e GOCACHE=/root/.cache/go-build -e REPO_URL -e GIT_BRANCH -e AWS_REGION -e AWS_DEFAULT_REGION -e AWS_ACCESS_KEY_ID -e AWS_SECRET_ACCESS_KEY -e AWS_SESSION_TOKEN -e AWS_SECURITY_TOKEN -e AWS_SESSION_EXPIRATION -e SKIP_SETUP -e SKIP_TEST -e SKIP_TEARDOWN $(BUILD_HARNESS_REPO):$(BUILD_HARNESS_VERSION) /bin/bash -c 'asdf install && go test -v -count=1 -timeout 2h -p 1 ./...'

.PHONY: docker-save-build-harness
docker-save-build-harness: ## Pulls the build harness docker image and saves it to a tarball
	@mkdir -p .cache/docker
	@docker pull $(BUILD_HARNESS_REPO):$(BUILD_HARNESS_VERSION)
	@docker save -o .cache/docker/build-harness.tar $(BUILD_HARNESS_REPO):$(BUILD_HARNESS_VERSION)

.PHONY: docker-load-build-harness
docker-load-build-harness: ## Loads the saved build harness docker image
	@docker load -i .cache/docker/build-harness.tar

.PHONY: run-pre-commit-hooks
run-pre-commit-hooks: ## Run all pre-commit hooks. Returns nonzero exit code if any hooks fail. Uses Docker for maximum compatibility
	@mkdir -p .cache/pre-commit
	@docker run --rm -v "${PWD}:/app" --workdir "/app" -e "PRE_COMMIT_HOME=/app/.cache/pre-commit" $(BUILD_HARNESS_REPO):$(BUILD_HARNESS_VERSION) pre-commit run -a --show-diff-on-failure

.PHONY: fix-cache-permissions
fix-cache-permissions: ## Fixes the permissions on the pre-commit cache
	@docker run --rm -v "${PWD}:/app" --workdir "/app" -e "PRE_COMMIT_HOME=/app/.cache/pre-commit" $(BUILD_HARNESS_REPO):$(BUILD_HARNESS_VERSION) chmod -R a+rx .cache

build:
	@mkdir -p build

build/zarf: | build ## Download the Linux flavor of Zarf to the build dir
	@echo "Downloading zarf"
	@curl -sL https://github.com/defenseunicorns/zarf/releases/download/$(ZARF_VERSION)/zarf_$(ZARF_VERSION)_Linux_amd64 -o build/zarf
	@chmod +x build/zarf

build/zarf-mac-intel: | build ## Download the Mac (Intel) flavor of Zarf to the build dir
	@echo "Downloading zarf-mac-intel"
	@curl -sL https://github.com/defenseunicorns/zarf/releases/download/$(ZARF_VERSION)/zarf_$(ZARF_VERSION)_Darwin_amd64 -o build/zarf-mac-intel
	@chmod +x build/zarf-mac-intel

build/zarf-init.sha256: | build ## Download the init package and create a small file with the sha256sum of the package so the Makefile can check whether it needs to be updated
	@echo "Downloading zarf-init-amd64-$(ZARF_VERSION).tar.zst"
	@curl -sL https://github.com/defenseunicorns/zarf/releases/download/$(ZARF_VERSION)/zarf-init-amd64-$(ZARF_VERSION).tar.zst -o build/zarf-init-amd64-$(ZARF_VERSION).tar.zst
	@echo "Creating shasum of the init package"
	@shasum -a 256 build/zarf-init-amd64-$(ZARF_VERSION).tar.zst | awk '{print $$1}' > build/zarf-init.sha256
