.PHONY : test build-web pack-web serve deploy build release build-all build-docker push-docker push-latest push-github release-notes check-env

test :
	go test ./lib

build-web :
	cp $$GOROOT/misc/wasm/wasm_exec.js web/assets/
	GOOS=js GOARCH=wasm go build -o web/assets/jd.wasm ./web/ui/main.go

pack-web : build-web
	go run web/pack/main.go

serve : pack-web
	go run -tags include_web main.go -port 8080

deploy : test build-web
	gsutil -m cp -r web/assets/* gs://play.jd-tool.io

build : test pack-web
	mkdir -p release
	CGO_ENABLED=0 go build -tags include_web -o release/jd main.go

release : check-env build-all push-github build-docker push-docker push-latest release-notes
	@echo
	@echo "Upload release/jd-* to Github as release $(JD_VERSION) with release notes above."
	@echo

build-all : test pack-web
	mkdir -p release
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -tags include_web -o release/jd-amd64-linux main.go
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -tags include_web -o release/jd-amd64-darwin main.go
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -tags include_web -o release/jd-amd64-windows main.go
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -tags include_web -o release/jd-arm64-linux main.go
	GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -tags include_web -o release/jd-arm64-darwin main.go
	GOOS=windows GOARCH=arm64 CGO_ENABLED=0 go build -tags include_web -o release/jd-arm64-windows main.go

build-docker : check-env test
	docker build -t josephburnett/jd:v$(JD_VERSION) .

push-docker : check-env build-docker
	docker push josephburnett/jd:v$(JD_VERSION)

push-latest : check-env build-docker
	docker tag josephburnett/jd:v$(JD_VERSION) josephburnett/jd:latest
	docker push josephburnett/jd:latest

push-github : check-env
	git diff --exit-code
	git tag v$(JD_VERSION) --force
	git push origin v$(JD_VERSION)

release-notes : check-env
	@echo
	@git log --oneline --no-decorate v$(JD_PREVIOUS_VERSION)..v$(JD_VERSION)

check-env :
ifndef JD_VERSION
	$(error Set version in main.go, commit and set JD_VERSION)
endif
ifndef JD_PREVIOUS_VERSION
	$(error Set JD_PREVIOUS_VERSION for release notes)
endif

