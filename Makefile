.PHONY : build test fuzz pack-web build-web serve release-build build-all build-docker release-push push-docker push-latest push-github deploy release-notes check-dirty check-version check-env find-issues

build : test pack-web
	mkdir -p release
	CGO_ENABLED=0 go build -tags include_web -o release/jd main.go

test :
	go test .
	go test ./lib
	cd v2 ; go test .

fuzz :
	go test ./lib -fuzz=FuzzJd

pack-web : build-web
	go run web/pack/main.go

build-web :
	cp $$GOROOT/misc/wasm/wasm_exec.js web/assets/
	GOOS=js GOARCH=wasm go build -o web/assets/jd.wasm ./web/ui/main.go

serve : pack-web
	go run -tags include_web main.go -port 8080

release-build : check-env check-version check-dirty build-all build-docker
	@echo
	@echo "If everything looks good, run 'make release-push'."
	@echo

build-all : test pack-web
	mkdir -p release
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -tags include_web -o release/jd-amd64-linux main.go
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -tags include_web -o release/jd-amd64-darwin main.go
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -tags include_web -o release/jd-amd64-windows.exe main.go
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -tags include_web -o release/jd-arm64-linux main.go
	GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -tags include_web -o release/jd-arm64-darwin main.go
	GOOS=windows GOARCH=arm64 CGO_ENABLED=0 go build -tags include_web -o release/jd-arm64-windows.exe main.go

build-docker : check-env test
	docker build -t josephburnett/jd:v$(JD_VERSION) .

release-push : check-env push-github push-docker push-latest deploy release-notes
	@echo
	@echo "Upload release/jd-* to Github as release $(JD_VERSION) with release notes above."
	@echo

push-docker : check-env
	docker push josephburnett/jd:v$(JD_VERSION)

push-latest : check-env
	docker tag josephburnett/jd:v$(JD_VERSION) josephburnett/jd:latest
	docker push josephburnett/jd:latest

push-github : check-env
	git diff --exit-code
	git tag v$(JD_VERSION) --force
	git push origin v$(JD_VERSION)

deploy : test build-web
	gsutil -m cp -r web/assets/* gs://play.jd-tool.io

release-notes : check-env
	@echo
	@git log --oneline --no-decorate v$(JD_PREVIOUS_VERSION)..v$(JD_VERSION)

check-dirty :
	git diff --quiet --exit-code

.ONESHELL:
check-version : check-env
	if ! grep -q $(JD_VERSION) main.go; then
		@echo "Set 'const version = $(JD_VERSION)' in main.go."
		false
	fi

check-env :
ifndef JD_VERSION
	$(error Set JD_VERSION)
endif
ifndef JD_PREVIOUS_VERSION
	$(error Set JD_PREVIOUS_VERSION for release notes)
endif

find-issues :
	-staticcheck ./...
	-goreportcard-cli -v ./...
