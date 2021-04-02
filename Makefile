.PHONY : test build-web pack-web serve deploy build release build-all build-docker push-docker check-env

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

release : check-env build-all build-docker push-docker
	echo "Upload release/jd-* to Github as release $$JD_VERSION"

build-all : test pack-web
	mkdir -p release
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -tags include_web -o release/jd-amd64-linux main.go
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -tags include_web -o release/jd-amd64-darwin main.go
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -tags include_web -o release/jd-amd64-windows main.go

build-docker : check-env test
	docker build -t josephburnett/jd:$$JD_VERSION .

push-docker : check-env
	docker push josephburnett/jd:$$JD_VERSION

check-env :
ifndef JD_VERSION
	echo "Tag the commit on a release branch and set JD_VERSION."
endif
