.goals = web
.PHONY : web

web :
	cp $$GOROOT/misc/wasm/wasm_exec.js web/assets/
	GOOS=js GOARCH=wasm go build -o web/assets/jd.wasm ./web/main.go
