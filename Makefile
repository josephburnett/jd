.goals = build-web pack-web serve test preflight deploy release
.PHONY : build-web pack-web serve test preflight deploy release

build-web :
	cp $$GOROOT/misc/wasm/wasm_exec.js web/assets/
	GOOS=js GOARCH=wasm go build -o web/assets/jd.wasm ./web/ui/main.go

pack-web : build-web
	go run web/pack/main.go

serve : pack-web
	go run -tags include_web main.go -port 8080

test :
	go test ./lib

preflight : test pack-web

deploy : preflight
	gsutil -m cp -r web/assets/* gs://play.jd-tool.io

release : preflight
	mkdir -p release
	CGO_ENABLED=0 go build -tags include_web -o release/jd main.go
