build:
	go build .

build-all:
	for arch in "amd64" "arm64" ; do \
		for os in "linux" "darwin" "windows" ; do \
			GOOS=$$os GOARCH=$$arch go build -o npkill-$$os-$$arch main.go ; \
		done \
	done


.PHONY: build build-all