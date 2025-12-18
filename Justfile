set dotenv-load

default:
  @just --list

pc: pre-commit

[group('qa')]
pre-commit: tidy generate lint

[group('qa')]
build:
  @go build -o ./build/terraform-provider-cloud ./main.go

[group('qa')]
lint:
  golangci-lint run --fix --build-tags it --timeout 5m

[group('qa')]
tidy:
  @go mod tidy

[group('test')]
tests: tests-unit tests-e2e tests-integration coverage

[group('test')]
coverage:
  @rm -rf coverage/coverage_merged.txt
  @head -n 1 coverage/coverage_unit.txt > coverage/coverage_merged.txt
  @tail -n +2 coverage/coverage_unit.txt | grep -Ev "generated|/pkg|tests/" >> coverage/coverage_merged.txt
  @tail -n +2 coverage/coverage_e2e.txt | grep -Ev "generated|/pkg|tests/" >> coverage/coverage_merged.txt
  @tail -n +2 coverage/coverage_integration.txt | grep -Ev "generated|/pkg|tests/" >> coverage/coverage_merged.txt
  @go tool cover -func=coverage/coverage_merged.txt

[group('test')]
generate:
  @go generate ./...

[group('test')]
tests-unit: 
  @mkdir -p ./coverage
  @go test -v -tags it ./internal/... -covermode=atomic -coverprofile=coverage/coverage_unit.txt -race -coverpkg=./internal/...

[group('test')]
tests-e2e tags="":
  @mkdir -p ./coverage
  @TF_ACC=1 go test -v ./tests/e2e/... -covermode=atomic -coverprofile=coverage/coverage_e2e.txt -race -coverpkg=./internal/...,./cmd/...

[group('test')]
tests-integration tags="":
  @mkdir -p ./coverage
  @TF_ACC=1 go test -v ./tests/integration/... -covermode=atomic -coverprofile=coverage/coverage_integration.txt -race -coverpkg=./internal/...,./cmd/...

[group('terraform')]
plan examples="install-verif":
  @go build -o ./build/terraform-provider-cloud ./main.go
  @cd examples/{{examples}} && terraform plan -generate-config-out=generated.tf

[group('terraform')]
apply examples="install-verif":
  @go build -o ./build/terraform-provider-cloud ./main.go
  @cd examples/{{examples}} && terraform apply -auto-approve

[group('terraform')]
destroy examples="install-verif":
  @go build -o ./build/terraform-provider-tf-cloud ./main.go
  @cd examples/{{examples}} && terraform destroy -auto-approve 

[group('releases')]
release-local: pc
  @goreleaser release --nightly --skip=publish --clean

[group('releases')]
release-ci: pc
  @goreleaser release --nightly --clean

[group('releases')]
release: pc
  @echo "$GPG_PRIVATE_KEY" | gpg --batch --import
  @echo "$GPG_FULL_FP:6:" | gpg --import-ownertrust -
  @goreleaser release --clean

[group('deployment')]
connect-dev:
  vcluster connect $USER --server=https://kube.$USER.formance.dev

generate-client:
  #!/usr/bin/env bash
  mkdir -p ./openapi
  rm -rf ./openapi/membership.yaml
  curl -o ./openapi/membership.yaml https://raw.githubusercontent.com/formancehq/membership-api/refs/heads/main/openapi.yaml -H "Authorization: Bearer $GITHUB_TOKEN"
  speakeasy generate sdk -s openapi/membership.yaml -o ./pkg/membership_client -l go