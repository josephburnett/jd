# Enforce strict toolchain usage - matches toolchain directive in go.mod
export GOTOOLCHAIN=go1.24.13

# Toolchain validation target
.PHONY : validate-toolchain
validate-toolchain :
	@echo "Validating Go toolchain versions..."
	@ROOT_TOOLCHAIN=$$(grep '^toolchain ' go.mod | awk '{print $$2}'); \
	V2_TOOLCHAIN=$$(grep '^toolchain ' v2/go.mod | awk '{print $$2}'); \
	if [ "$$ROOT_TOOLCHAIN" != "$$V2_TOOLCHAIN" ]; then \
		echo "Error: Toolchain mismatch between go.mod files:"; \
		echo "  Root go.mod: $$ROOT_TOOLCHAIN"; \
		echo "  v2/go.mod: $$V2_TOOLCHAIN"; \
		exit 1; \
	fi; \
	if [ "$$ROOT_TOOLCHAIN" != "$(GOTOOLCHAIN)" ]; then \
		echo "Error: Makefile GOTOOLCHAIN does not match go.mod toolchain:"; \
		echo "  Makefile GOTOOLCHAIN: $(GOTOOLCHAIN)"; \
		echo "  go.mod toolchain: $$ROOT_TOOLCHAIN"; \
		echo "  Please update GOTOOLCHAIN in Makefile to match go.mod"; \
		exit 1; \
	fi; \
	GOMOD_VERSION=$$(grep '^go ' go.mod | awk '{print $$2}' | cut -d. -f1,2); \
	DOCKER_VERSION=$$(grep '^FROM golang:' Dockerfile | head -1 | sed 's/FROM golang:\([0-9.]*\).*/\1/'); \
	if [ "$$GOMOD_VERSION" != "$$DOCKER_VERSION" ]; then \
		echo "Error: Dockerfile Go version does not match go.mod:"; \
		echo "  Dockerfile: golang:$$DOCKER_VERSION"; \
		echo "  go.mod go directive: $$GOMOD_VERSION"; \
		echo "  Please update FROM golang: in Dockerfile to match go.mod"; \
		exit 1; \
	fi; \
	echo "✓ Toolchain validation passed: $$ROOT_TOOLCHAIN"

.PHONY : build
build : test pack-web validate-toolchain
	mkdir -p release
	cd v2/jd ; CGO_ENABLED=0 go build -tags include_web -o ../../release/jd main.go

.PHONY : test
test : validate-toolchain
	go test .
	go test ./lib
	cd v2 ; go test -run '^Test' .
	cd v2 ; go test ./jd

.PHONY : fuzz
fuzz : validate-toolchain
	cd v2 ; go test . -fuzz=FuzzJd -fuzztime=5m

.PHONY : fuzz-indef
fuzz-indef: validate-toolchain
	cd v2 ; go test . -fuzz=FuzzJd

.PHONY : cover
cover : validate-toolchain
	cd v2 ; go test -coverprofile=coverage.out . ./jd
	cd v2 ; go run internal/coverfilter/main.go coverage.out > coverage_filtered.out
	cd v2 ; go tool cover -func=coverage_filtered.out | tail -1
	@echo "For HTML report: cd v2 && go tool cover -html=coverage.out"
	@echo "Checking non-trivial code for 100% coverage..."
	@cd v2 ; grep -v '^\s*#' cover_exclude.txt | grep -v '^\s*$$' > coverage_exclude_patterns.tmp ; \
	lines=$$(go tool cover -func=coverage_filtered.out \
		| grep -v '^total:' \
		| grep -v '100.0%' \
		| grep -vf coverage_exclude_patterns.tmp) ; \
	rm -f coverage_exclude_patterns.tmp ; \
	if [ -n "$$lines" ]; then \
		echo "$$lines" ; \
		echo "FAIL: non-trivial functions below 100% coverage" ; \
		exit 1 ; \
	else \
		echo "OK: all non-trivial functions at 100% coverage" ; \
	fi

.PHONY : spec-test
spec-test : validate-toolchain
	cd v2 ; go build -o ../spec/test/jd ./jd
	cd spec/test ; go build -o test-runner .
	cd spec/test ; ./test-runner -verbose -binary ./jd
	rm -f spec/test/jd spec/test/test-runner

.PHONY : go-fmt
go-fmt :
	go fmt ./...
	cd v2 ; go fmt ./...

.PHONY : check-fmt
check-fmt : validate-toolchain
	@unformatted=$$(gofmt -l . v2); \
	if [ -n "$$unformatted" ]; then \
		echo "Error: unformatted Go files:"; \
		echo "$$unformatted"; \
		echo "Run 'make go-fmt' to fix."; \
		exit 1; \
	fi; \
	echo "✓ All Go files formatted"

.PHONY : vet
vet : validate-toolchain
	go vet . ./lib
	cd v2 ; go vet $$(go list ./... 2>/dev/null | grep -v internal/web/ui)

.PHONY : pack-web
pack-web : build-web validate-toolchain
	cd v2 ; go run internal/web/pack/main.go

.PHONY : build-web
build-web : validate-toolchain
	cp "$$(go env GOROOT)/lib/wasm/wasm_exec.js" v2/internal/web/assets/wasm_exec.js
	cd v2 ; GOOS=js GOARCH=wasm go build -o internal/web/assets/jd.wasm ./internal/web/ui

.PHONY : serve
serve : pack-web validate-toolchain
	cd v2 ; go run -tags include_web jd/main.go -port 8080

.PHONY : release-build
release-build : check-env check-version check-dirty build-all build-docker
	@echo
	@echo "If everything looks good, run 'make release-push'."
	@echo

.PHONY : build-all
build-all : test pack-web validate-toolchain
	mkdir -p release
	cd v2/jd ; GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -tags include_web -o ../../release/jd-amd64-linux main.go
	cd v2/jd ; GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -tags include_web -o ../../release/jd-amd64-darwin main.go
	cd v2/jd ; GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -tags include_web -o ../../release/jd-amd64-windows.exe main.go
	cd v2/jd ; GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -tags include_web -o ../../release/jd-arm64-linux main.go
	cd v2/jd ; GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -tags include_web -o ../../release/jd-arm64-darwin main.go
	cd v2/jd ; GOOS=windows GOARCH=arm64 CGO_ENABLED=0 go build -tags include_web -o ../../release/jd-arm64-windows.exe main.go
	cd v2/jd ; GOOS=linux GOARCH=riscv64 CGO_ENABLED=0 go build -tags include_web -o ../../release/jd-riscv64-linux main.go

.PHONY : build-docker
build-docker : check-env test
	docker build -t josephburnett/jd:v$(JD_VERSION) .

.PHONY : release-push
release-push : check-env push-github push-docker push-latest deploy release-notes
	@echo
	@echo "Upload release/jd-* to Github as release $(JD_VERSION) with release notes above."
	@echo

.PHONY : push-docker
push-docker : check-env
	docker push josephburnett/jd:v$(JD_VERSION)

.PHONY : push-latest
push-latest : check-env
	docker tag josephburnett/jd:v$(JD_VERSION) josephburnett/jd:latest
	docker push josephburnett/jd:latest

.PHONY : push-github
push-github : check-env
	git diff --exit-code
	git tag v$(JD_VERSION) --force
	git push origin v$(JD_VERSION)

.PHONY : deploy
deploy : test build-web
	gsutil -m cp -r v2/internal/web/assets/* gs://play.jd-tool.io

.PHONY : release-notes
release-notes : check-env
	@echo
	@git log --oneline --no-decorate v$(JD_PREVIOUS_VERSION)..v$(JD_VERSION)

.PHONY : check-dirty
check-dirty : tidy
	git diff --quiet --exit-code

.PHONY : tidy
tidy : validate-toolchain
	cd v2 ; go mod tidy
	go mod tidy

.PHONY : check-version
check-version : check-env
	@if ! grep -q $(JD_VERSION) v2/jd/main.go; then                   \
		echo "Set 'const version = $(JD_VERSION)' in main.go." ;  \
		false                                                   ; \
	fi
	@if ! grep -q v$(JD_VERSION) action.yml; then                                  \
		echo "Set 'docker://josephburnett/jd:v$(JD_VERSION)' in action.yml." ; \
		false                                                   ;              \
	fi

.PHONY : check-env
check-env :
ifndef JD_VERSION
	$(error Set JD_VERSION)
endif
ifndef JD_PREVIOUS_VERSION
	$(error Set JD_PREVIOUS_VERSION for release notes)
endif

.PHONY : benchmark
benchmark : validate-toolchain
	@echo "Running performance baseline benchmarks..."
	@mkdir -p benchmarks
	@timestamp=$$(date +%Y%m%d_%H%M%S); \
	cd v2 && go test -run=^$$ -bench=^Benchmark -benchmem -count=1 -timeout=3m -benchtime=200ms | tee ../benchmarks/baseline_$$timestamp.txt; \
	echo "Results saved to benchmarks/baseline_$$timestamp.txt"

.PHONY : find-issues
find-issues :
	-staticcheck ./...
	-goreportcard-cli -v ./...
