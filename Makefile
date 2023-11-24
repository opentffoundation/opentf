export PATH := $(abspath bin/):${PATH}

# Dependency versions
LICENSEI_VERSION = 0.9.0

# generate runs `go generate` to build the dynamically generated
# source files, except the protobuf stubs which are built instead with
# "make protobuf".
.PHONY: generate
generate:
	go generate ./...

# We separate the protobuf generation because most development tasks on
# OpenTofu do not involve changing protobuf files and protoc is not a
# go-gettable dependency and so getting it installed can be inconvenient.
#
# If you are working on changes to protobuf interfaces, run this Makefile
# target to be sure to regenerate all of the protobuf stubs using the expected
# versions of protoc and the protoc Go plugins.
.PHONY: protobuf
protobuf:
	go run ./tools/protobuf-compile .

.PHONY: fmtcheck
fmtcheck:
	"$(CURDIR)/scripts/gofmtcheck.sh"

.PHONY: importscheck
importscheck:
	"$(CURDIR)/scripts/goimportscheck.sh"

.PHONY: staticcheck
staticcheck:
	"$(CURDIR)/scripts/staticcheck.sh"

.PHONY: exhaustive
exhaustive:
	"$(CURDIR)/scripts/exhaustive.sh"

# Run license check
.PHONY: license-check
license-check:
	go mod vendor
	licensei cache --debug
	licensei check --debug
	licensei header --debug
	rm -rf vendor/
	git diff --exit-code

# Install dependencies
deps: bin/licensei
deps:

bin/licensei: bin/licensei-${LICENSEI_VERSION}
	@ln -sf licensei-${LICENSEI_VERSION} bin/licensei
bin/licensei-${LICENSEI_VERSION}:
	@mkdir -p bin
	curl -sfL https://git.io/licensei | bash -s v${LICENSEI_VERSION}
	@mv bin/licensei $@

# disallow any parallelism (-j) for Make. This is necessary since some
# commands during the build process create temporary files that collide
# under parallel conditions.
.NOTPARALLEL:

# Integration tests
#

.PHONY:
list-integration-tests: ## Lists tests.
	@ grep -h -E '^(test|integration)-.+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[1m%-30s\033[0m %s\n", $$1, $$2}'

# integration test with s3 as backend
.PHONY: test-s3

define infoTestS3
Test requires:
* AWS Credentials to be configured
  - https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-files.html
  - https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-envvars.html
* IAM Permissions in us-west-2
  - S3 CRUD operations on buckets which will follow the pattern tofu-test-*
  - DynamoDB CRUD operations on a Table named dynamoTable

endef

test-s3: ## Runs tests with s3 bucket as the backend.
	@ $(info $(infoTestS3))
	@ TF_S3_TEST=1 go test ./internal/backend/remote-state/s3/...

# integration test with postgres as backend
.PHONY: test-pg test-pg-clean

PG_PORT := 5432

define infoTestPg
Test requires:
* Docker: https://docs.docker.com/engine/install/
* Port: $(PG_PORT)

endef

test-pg: test-pg-clean ## Runs tests with local Postgres instance as the backend.
	@ $(info $(infoTestPg))
	@ docker run --rm -d --name tofu-pg \
        -p $(PG_PORT):5432 \
        -e POSTGRES_PASSWORD=tofu \
        -e POSTGRES_USER=tofu \
        postgres:14-alpine3.17
	@ docker exec tofu-pg /bin/bash -c 'until psql -U tofu -c "\q"; do >&2 echo "db is getting ready, waiting"; sleep 1; done'
	@ DATABASE_URL="postgres://tofu:tofu@localhost:$(PG_PORT)/tofu?sslmode=disable" \
 		TF_PG_TEST=1 TF_ACC=1 go test ./internal/backend/remote-state/pg/...

test-pg-clean: ## Cleans environment after `test-pg`.
	@ echo "Cleans after test-pg"
	@ docker rm -f tofu-pg 1> /dev/null

.PHONY:
integration-tests: test-pg integration-tests-clean ## Runs all integration tests test.

.PHONY:
integration-tests-clean: test-pg-clean ## Cleans environment after all integration tests.
