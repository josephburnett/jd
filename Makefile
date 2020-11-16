.goals = build-web deploy serve
.PHONY : build-web deploy serve

build-web :
	cp $$GOROOT/misc/wasm/wasm_exec.js web/assets/
	GOOS=js GOARCH=wasm go build -o web/assets/jd.wasm ./web/main.go

deploy : build-web
	gsutil -m cp -r web/assets/* gs://play.jd-tool.io

serve : build-web
	go run main.go -port 8080
