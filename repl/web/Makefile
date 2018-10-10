build:
	GOARCH=wasm GOOS=js go build -o gpython.wasm
	gopherjs build -m -o gpython.js

serve:	build
	go run serve.go

upload:	build
	rclone -P copy --include="*.{wasm,js,css,html}" . gpythonweb:
	@echo See https://gpython.org/
