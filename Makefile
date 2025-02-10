.PHONY : build
build : test pack-web
	mkdir -p release
	CGO_ENABLED=0 go build -tags include_web -o release/jd main.go

.PHONY : test
test :
	go test .
	go test ./lib
	cd v2 ; go test .

.PHONY : fuzz
fuzz :
	go test ./lib -fuzz=FuzzJd -fuzztime=10s
	cd v2 ; go test . -fuzz=FuzzJd -fuzztime=10s

.PHONY : pack-web
pack-web : build-web
	go run web/pack/main.go

.PHONY : build-web
build-web : check-goroot
	cp $$GOROOT/misc/wasm/wasm_exec.js web/assets/
	GOOS=js GOARCH=wasm go build -o web/assets/jd.wasm ./web/ui/main.go

.PHONY : serve
serve : pack-web
	go run -tags include_web main.go -port 8080

.PHONY : release-build
release-build : check-env check-version check-dirty build-all build-docker
	@echo
	@echo "If everything looks good, run 'make release-push'."
	@echo

.PHONY : build-all
build-all : test pack-web
	mkdir -p release
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -tags include_web -o release/jd-amd64-linux main.go
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -tags include_web -o release/jd-amd64-darwin main.go
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -tags include_web -o release/jd-amd64-windows.exe main.go
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -tags include_web -o release/jd-arm64-linux main.go
	GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -tags include_web -o release/jd-arm64-darwin main.go
	GOOS=windows GOARCH=arm64 CGO_ENABLED=0 go build -tags include_web -o release/jd-arm64-windows.exe main.go

.PHONY : build-docker
build-docker : check-env test
	docker build -t josephburnett/jd:v$(JD_VERSION) .

.PHONY : release-push
release-push : check-env push-github push-docker release-notes
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
	gsutil -m cp -r web/assets/* gs://play.jd-tool.io

.PHONY : release-notes
release-notes : check-env
	@echo
	@git log --oneline --no-decorate v$(JD_PREVIOUS_VERSION)..v$(JD_VERSION)

.PHONY : check-dirty
check-dirty : vendor
	git diff --quiet --exit-code

.PHONY : vendor
vendor :
	cd v2 ; go mod tidy
	cd v2 ; go mod vendor
	go mod tidy
	go mod vendor

.PHONY : check-version
check-version : check-env
	@if ! grep -q $(JD_VERSION) main.go; then                          \
		echo "Set 'const version = $(JD_VERSION)' in main.go." ; \
		false                                                   ; \
	fi

.PHONY : check-env
check-env :
ifndef JD_VERSION
	$(error Set JD_VERSION)
endif
ifndef JD_PREVIOUS_VERSION
	$(error Set JD_PREVIOUS_VERSION for release notes)
endif

.PHONY : check-goroot
check-goroot:
ifndef GOROOT
	$(error Set GOROOT in order to file wasm_exec.js)
endif

.PHONY : find-issues
find-issues :
	-staticcheck ./...
	-goreportcard-cli -v ./...
